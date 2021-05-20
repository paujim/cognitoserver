package services

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/paujim/cognitoserver/server/pkg/entities"
	log "github.com/sirupsen/logrus"
)

var (
	ErrorInvalidInputParameters = errors.New("Missing input parameters")
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

type cognitoHandler struct {
	appClientID *string
	userPoolID  *string
	cognitoAPI  cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

func NewCognitoHandler(appClientID, userPoolID string, client cognitoidentityprovideriface.CognitoIdentityProviderAPI) entities.UserTokenHandler {
	return &cognitoHandler{
		appClientID: aws.String(appClientID),
		userPoolID:  aws.String(userPoolID),
		cognitoAPI:  client,
	}
}

func (c *cognitoHandler) GetTokens(username, password *string) (accessToken, refreshToken *string, err error) {

	if username == nil || password == nil {
		err = ErrorInvalidInputParameters
		return
	}

	log.Info("Getting access token")
	params := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: c.appClientID,
		AuthParameters: map[string]*string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	}
	req, resp := c.cognitoAPI.InitiateAuthRequest(params)
	err = req.Send()
	if err != nil {
		return
	}
	log.Info(resp.GoString())

	// Ok
	if resp.ChallengeName == nil {
		accessToken = resp.AuthenticationResult.AccessToken
		refreshToken = resp.AuthenticationResult.RefreshToken
		return
	}
	// NEW_PASSWORD_REQUIRED Challenge
	if *resp.ChallengeName == "NEW_PASSWORD_REQUIRED" {
		return c.responseToNewPassword(resp.Session, username, password)
	}
	// Others
	err = errors.New("Unable to respond: " + *resp.ChallengeName)
	return
}

func (c *cognitoHandler) responseToNewPassword(session, username, password *string) (accessToken, refreshToken *string, err error) {
	log.Infoln("New password required. Responding chanllenge with old password")
	params := &cognitoidentityprovider.RespondToAuthChallengeInput{
		Session:       session,
		ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
		ClientId:      c.appClientID,
		ChallengeResponses: map[string]*string{
			"USERNAME":     username,
			"NEW_PASSWORD": password,
		},
	}
	req, resp := c.cognitoAPI.RespondToAuthChallengeRequest(params)
	err = req.Send()
	if err != nil {
		return
	}

	if resp.AuthenticationResult == nil {
		err = errors.New("Unable to get AccessToken")
		return
	}

	log.Info(resp.GoString())
	accessToken = resp.AuthenticationResult.AccessToken
	refreshToken = resp.AuthenticationResult.RefreshToken
	return
}

func (c *cognitoHandler) RefreshAccessToken(token *string) (accessToken, refreshToken *string, err error) {

	if token == nil {
		err = ErrorInvalidInputParameters
		return
	}

	log.Info("Refreshing token")
	params := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("REFRESH_TOKEN_AUTH"),
		ClientId: c.appClientID,
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": token,
		},
	}
	req, resp := c.cognitoAPI.InitiateAuthRequest(params)
	err = req.Send()
	if err != nil {
		return
	}
	log.Info(resp.GoString())
	if resp.AuthenticationResult == nil {
		err = errors.New("Unable to get AccessToken")
		return
	}
	accessToken = resp.AuthenticationResult.AccessToken
	refreshToken = token
	return
}

func (c *cognitoHandler) RegisterUser(username, password *string) (sub *string, err error) {
	log.Info("Registering new user")
	params := &cognitoidentityprovider.SignUpInput{
		ClientId: c.appClientID,
		Password: password,
		Username: username,
	}
	req, resp := c.cognitoAPI.SignUpRequest(params)
	err = req.Send()
	if err != nil {
		return
	}
	log.Info(resp.GoString())
	sub = resp.UserSub
	return
}

func (c *cognitoHandler) ListUsers() (users []entities.UserModel, err error) {
	log.Info("Getting all users")
	params := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: c.userPoolID,
	}

	req, resp := c.cognitoAPI.ListUsersRequest(params)
	err = req.Send()
	if err != nil {
		return
	}
	log.Info(resp.GoString())

	users = []entities.UserModel{}
	for _, user := range resp.Users {
		users = append(users, entities.UserModel{
			Username: user.Username,
			Status:   user.UserStatus,
			Enabled:  user.Enabled,
			Created:  user.UserCreateDate,
		})
	}
	return
}
