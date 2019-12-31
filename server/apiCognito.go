package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"time"
)

//CognitoParam ...
type CognitoParam struct {
	appClientID *string
	userPoolID  *string
	client      cognitoidentityprovideriface.CognitoIdentityProviderAPI
	logger      ILogger
}

//NewCognitoParam ...
func NewCognitoParam(region, appClientID, userPoolID string, logger ILogger) *CognitoParam {

	mySession := session.Must(session.NewSession())
	cognito := &CognitoParam{
		appClientID: aws.String(appClientID),
		userPoolID:  aws.String(userPoolID),
		client:      cognitoidentityprovider.New(mySession, aws.NewConfig().WithRegion(region)),
		logger:      logger,
	}
	return cognito
}

//GetTokens ...
func (c *CognitoParam) GetTokens(username, password *string) (accessToken, refreshToken *string, err error) {

	c.logger.Logln("Getting access token")
	params := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: c.appClientID,
		AuthParameters: map[string]*string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	}
	req, resp := c.client.InitiateAuthRequest(params)
	err = req.Send()
	if err != nil {
		return
	}
	c.logger.Logln("\n" + resp.GoString())

	if resp.ChallengeName != nil && *resp.ChallengeName == "NEW_PASSWORD_REQUIRED" {
		return c.ResponseToNewPassword(resp.Session, username, password)
	}

	accessToken = resp.AuthenticationResult.AccessToken
	refreshToken = resp.AuthenticationResult.RefreshToken
	return
}

//ResponseToNewPassword ...
func (c *CognitoParam) ResponseToNewPassword(session, username, password *string) (accessToken, refreshToken *string, err error) {
	c.logger.Logln("New password required. Responding chanllenge with old password")
	params := &cognitoidentityprovider.RespondToAuthChallengeInput{
		Session:       session,
		ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
		ClientId:      c.appClientID,
		ChallengeResponses: map[string]*string{
			"USERNAME":     username,
			"NEW_PASSWORD": password,
		},
	}
	req, resp := c.client.RespondToAuthChallengeRequest(params)
	err = req.Send()
	if err != nil {
		return
	}

	c.logger.Logln("\n" + resp.GoString())
	accessToken = resp.AuthenticationResult.AccessToken
	refreshToken = resp.AuthenticationResult.RefreshToken
	return
}

//RefreshAccessToken ...
func (c *CognitoParam) RefreshAccessToken(token *string) (accessToken, refreshToken *string, err error) {

	c.logger.Logln("Refreshing token")
	params := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("REFRESH_TOKEN_AUTH"),
		ClientId: c.appClientID,
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": token,
		},
	}
	req, resp := c.client.InitiateAuthRequest(params)
	err = req.Send()
	if err != nil {
		return
	}
	c.logger.Logln("\n" + resp.GoString())

	accessToken = resp.AuthenticationResult.AccessToken
	refreshToken = token
	return
}

//RegisterUser ...
func (c *CognitoParam) RegisterUser(username, password *string) (sub *string, err error) {
	c.logger.Logln("Registering new user")
	params := &cognitoidentityprovider.SignUpInput{
		ClientId: c.appClientID,
		Password: password,
		Username: username,
	}
	req, resp := c.client.SignUpRequest(params)
	err = req.Send()
	if err != nil {
		return
	}
	c.logger.Logln("\n" + resp.GoString())
	sub = resp.UserSub
	return
}

//UserModel ...
type UserModel struct {
	Username *string    `json:"username"`
	Status   *string    `json:"status"`
	Enabled  *bool      `json:"enabled"`
	Created  *time.Time `json:"created"`
}

//ListUsers ...
func (c *CognitoParam) ListUsers() (users []UserModel, err error) {
	c.logger.Logln("Getting all users")
	params := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: c.userPoolID,
	}

	req, resp := c.client.ListUsersRequest(params)
	err = req.Send()
	if err != nil {
		return
	}
	c.logger.Logln("\n" + resp.GoString())

	users = []UserModel{}
	for _, user := range resp.Users {
		users = append(users, UserModel{
			Username: user.Username,
			Status:   user.UserStatus,
			Enabled:  user.Enabled,
			Created:  user.UserCreateDate,
		})
	}

	return
}