package controllers

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/paujim/cognitoserver/server/pkg/entities"
)

type auth struct {
	userPoolRegion string
	userPoolID     string
	service        entities.TokenHandler
}

func NewAuth(region, userPoolID string, service entities.TokenHandler) *auth {
	return &auth{
		userPoolRegion: region,
		userPoolID:     userPoolID,
		service:        service,
	}
}

func (a *auth) RegisterAuthRoutes(router *gin.RouterGroup) {
	router.POST("/token", a.getAccessToken)
}

func (a *auth) AuthMiddleware() gin.HandlerFunc {
	//Download and store the JSON Web Key (JWK) for your user pool.
	jwkURL := fmt.Sprintf("https://cognito-idp.%v.amazonaws.com/%v/.well-known/jwks.json", a.userPoolRegion, a.userPoolID)
	log.Println(jwkURL)
	jwk := a.getJWK(jwkURL)

	return func(c *gin.Context) {
		tokenString, ok := a.getBearer(c.Request.Header["Authorization"])
		if !ok {
			// Authorization Bearer Header is missing
			c.AbortWithStatusJSON(401, gin.H{"error": "missing_authorization_header"})
			return
		}

		token, err := a.validateToken(tokenString, jwk)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid_token"})
		} else {
			// All Good :)
			c.Set("token", token)
			c.Next()
		}
	}
}

type jwkKey struct {
	Alg string
	E   string
	Kid string
	Kty string
	N   string
	Use string
}

func (a *auth) getJWK(jwkURL string) map[string]jwkKey {
	type JWK struct {
		Keys []jwkKey
	}
	getJSON := func(url string, target interface{}) error {
		var myClient = &http.Client{Timeout: 10 * time.Second}
		r, err := myClient.Get(url)
		if err != nil {
			return err
		}
		defer r.Body.Close()

		return json.NewDecoder(r.Body).Decode(target)
	}

	jwk := &JWK{}

	getJSON(jwkURL, jwk)

	jwkMap := make(map[string]jwkKey)
	for _, jwk := range jwk.Keys {
		jwkMap[jwk.Kid] = jwk
	}
	return jwkMap
}

func (a *auth) getBearer(auth []string) (jwt string, ok bool) {
	for _, v := range auth {
		ret := strings.Split(v, " ")
		if len(ret) > 1 && ret[0] == "Bearer" {
			return ret[1], true
		}
	}
	return "", false
}

func (a *auth) validateToken(tokenStr string, jwk map[string]jwkKey) (*jwt.Token, error) {

	//Decode the token string into JWT format.
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		// cognito user pool : RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Get the kid from the JWT token header and retrieve the corresponding JSON Web Key that was stored
		if kid, ok := token.Header["kid"]; ok {
			if kidStr, ok := kid.(string); ok {
				key := jwk[kidStr]
				// Verify the signature of the decoded JWT token.
				rsaPublicKey := publicKey(key.E, key.N)
				return rsaPublicKey, nil
			}
		}

		// rsa public key
		return "", nil
	})

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	_, ok := claims["iss"]
	if !ok {
		return nil, fmt.Errorf("token does not contain issuer")
	}

	// Check the iss claim. It should match your user pool.
	issShoudBe := fmt.Sprintf("https://cognito-idp.%v.amazonaws.com/%v", a.userPoolRegion, a.userPoolID)
	err = a.validateClaimItem("iss", []string{issShoudBe}, claims)
	if err != nil {
		return nil, err
	}

	// Check the token_use claim.
	err = a.validateTokenUse(claims)
	if err != nil {
		return nil, err
	}

	// Check the exp claim and make sure the token is not expired.
	err = a.validateExpired(claims)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func publicKey(rawE, rawN string) *rsa.PublicKey {
	decodedE, err := base64.RawURLEncoding.DecodeString(rawE)
	if err != nil {
		panic(err)
	}
	if len(decodedE) < 4 {
		ndata := make([]byte, 4)
		copy(ndata[4-len(decodedE):], decodedE)
		decodedE = ndata
	}
	pubKey := &rsa.PublicKey{
		N: &big.Int{},
		E: int(binary.BigEndian.Uint32(decodedE[:])),
	}
	decodedN, err := base64.RawURLEncoding.DecodeString(rawN)
	if err != nil {
		panic(err)
	}
	pubKey.N.SetBytes(decodedN)
	return pubKey
}

func (a *auth) validateTokenUse(claims jwt.MapClaims) error {
	if tokenUse, ok := claims["token_use"]; ok {
		if tokenUseStr, ok := tokenUse.(string); ok {
			if tokenUseStr == "access" {
				return nil
			}
		}
	}
	return errors.New("token_use should be access")
}

func (a *auth) validateClaimItem(key string, keyShouldBe []string, claims jwt.MapClaims) error {
	if val, ok := claims[key]; ok {
		if valStr, ok := val.(string); ok {
			for _, shouldbe := range keyShouldBe {
				if valStr == shouldbe {
					return nil
				}
			}
		}
	}
	return fmt.Errorf("%v does not match any of valid values: %v", key, keyShouldBe)
}

func (a *auth) validateExpired(claims jwt.MapClaims) error {
	if tokenExp, ok := claims["exp"]; ok {
		if exp, ok := tokenExp.(float64); ok {
			now := time.Now().Unix()
			if int64(exp) > now {
				return nil
			}
		}
		return errors.New("cannot parse token exp")
	}
	return errors.New("token is expired")
}

func (a *auth) getAccessToken(c *gin.Context) {
	var request entities.TokenRequest
	c.ShouldBind(&request)

	var accessToken, refreshToken *string
	var err error

	if request.RefreshToken == nil {
		accessToken, refreshToken, err = a.service.GetTokens(request.Username, request.Password)
	} else {
		accessToken, refreshToken, err = a.service.RefreshAccessToken(request.RefreshToken)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token_type":    "bearer",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	return
}
