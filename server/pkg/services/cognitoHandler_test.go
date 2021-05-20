package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

// mocks
type mockedCognitoClient struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	initiateAuthRequest           *request.Request
	initiateAuthOutput            *cognitoidentityprovider.InitiateAuthOutput
	respondToAuthChallengeRequest *request.Request
	respondToAuthChallengeOutput  *cognitoidentityprovider.RespondToAuthChallengeOutput
	listUsersRequest              *request.Request
	listUsersRequestOutput        *cognitoidentityprovider.ListUsersOutput
}

func (m *mockedCognitoClient) InitiateAuthRequest(*cognitoidentityprovider.InitiateAuthInput) (*request.Request, *cognitoidentityprovider.InitiateAuthOutput) {
	return m.initiateAuthRequest, m.initiateAuthOutput
}
func (m *mockedCognitoClient) RespondToAuthChallengeRequest(*cognitoidentityprovider.RespondToAuthChallengeInput) (*request.Request, *cognitoidentityprovider.RespondToAuthChallengeOutput) {
	return m.respondToAuthChallengeRequest, m.respondToAuthChallengeOutput
}
func (m *mockedCognitoClient) ListUsersRequest(*cognitoidentityprovider.ListUsersInput) (*request.Request, *cognitoidentityprovider.ListUsersOutput) {
	return m.listUsersRequest, m.listUsersRequestOutput
}

func TestGetTokens(t *testing.T) {
	authResult := &cognitoidentityprovider.AuthenticationResultType{
		AccessToken:  aws.String("ACCESS_TOKEN"),
		RefreshToken: aws.String("REFRESH_TOKEN"),
	}
	expectedError := errors.New("Something went wrong")

	t.Run("Missing parameters on GetTokens", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{},
		)
		_, _, err := cp.GetTokens(nil, nil)
		if err != ErrorInvalidInputParameters {
			t.Errorf("Expected error when nil parameters")
		}
	})
	t.Run("Successfull GetTokens", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
					AuthenticationResult: authResult,
				},
			},
		)
		accessToken, refreshToken, err := cp.GetTokens(aws.String("username"), aws.String("password"))
		if err != nil {
			t.Errorf(err.Error())
		}
		if accessToken != nil && *accessToken != "ACCESS_TOKEN" {
			t.Errorf("Access token does not match the expected value")
		}
		if refreshToken != nil && *refreshToken != "REFRESH_TOKEN" {
			t.Errorf("The refresh token does not match the expected value")
		}
	})
	t.Run("Fail GetTokens", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{Error: expectedError},
				initiateAuthOutput:  nil,
			},
		)
		_, _, err := cp.GetTokens(aws.String("username"), aws.String("password"))

		if err != expectedError {
			t.Errorf("Expected error")
		}
	})
	t.Run("Successfull GetTokens with NEW_PASSWORD_REQUIRED", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
					ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
				},
				respondToAuthChallengeRequest: &request.Request{},
				respondToAuthChallengeOutput: &cognitoidentityprovider.RespondToAuthChallengeOutput{
					AuthenticationResult: authResult,
				},
			},
		)
		accessToken, refreshToken, err := cp.GetTokens(aws.String("username"), aws.String("password"))
		if err != nil {
			t.Errorf(err.Error())
		}
		if accessToken != nil && *accessToken != "ACCESS_TOKEN" {
			t.Errorf("Access token does not match the expected value")
		}
		if refreshToken != nil && *refreshToken != "REFRESH_TOKEN" {
			t.Errorf("The refresh token does not match the expected value")
		}
	})
	t.Run("Fail GetTokens with NEW_PASSWORD_REQUIRED", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
					ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
				},
				respondToAuthChallengeRequest: &request.Request{Error: expectedError},
			},
		)
		_, _, err := cp.GetTokens(aws.String("username"), aws.String("password"))
		if err != expectedError {
			t.Errorf("Ëxpected error")
		}
	})
	t.Run("Fail GetTokens with OTHER challenge", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
					ChallengeName:        aws.String("OTHER"),
					AuthenticationResult: authResult,
				},
			},
		)
		_, _, err := cp.GetTokens(aws.String("username"), aws.String("password"))
		if err == nil {
			t.Errorf("Error expected")
		}
	})
}

func TestRefreshAccessToken(t *testing.T) {
	authResult := &cognitoidentityprovider.AuthenticationResultType{
		AccessToken: aws.String("ACCESS_TOKEN"),
	}
	expectedError := errors.New("Something went wrong")
	t.Run("Missing parameters on RefreshAccessToken", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{},
		)
		_, _, err := cp.RefreshAccessToken(nil)
		if err != ErrorInvalidInputParameters {
			t.Errorf("Expected error when nil parameters")
		}
	})
	t.Run("Successfull RefreshAccessToken", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
					AuthenticationResult: authResult,
				},
			},
		)
		accessToken, refreshToken, err := cp.RefreshAccessToken(aws.String("refresh_token"))
		if err != nil {
			t.Errorf(err.Error())
		}
		if accessToken != nil && *accessToken != "ACCESS_TOKEN" {
			t.Errorf("Access token does not match the expected value")
		}
		if refreshToken != nil && *refreshToken != "refresh_token" {
			t.Errorf("The refresh token does not match the expected value")
		}
	})
	t.Run("Fail RefreshAccessToken with OTHER challenge", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
					ChallengeName: aws.String("OTHER"),
				},
			},
		)
		_, _, err := cp.RefreshAccessToken(aws.String("refresh_token"))
		if err == nil {
			t.Errorf("Error expected")
		}
	})
	t.Run("Fail RefreshAccessToken", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{Error: expectedError},
			},
		)
		_, _, err := cp.RefreshAccessToken(aws.String("refresh_token"))

		if err != expectedError {
			t.Errorf("Expected error")
		}
	})
}

func TestListUsers(t *testing.T) {
	expectedError := errors.New("Something went wrong")
	t.Run("Successfull ListUsers with three users", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{
				listUsersRequest: &request.Request{},
				listUsersRequestOutput: &cognitoidentityprovider.ListUsersOutput{
					Users: []*cognitoidentityprovider.UserType{
						{Username: aws.String("username_1")},
						{Username: aws.String("username_2")},
						{Username: aws.String("username_3")},
					},
				},
			},
		)
		users, err := cp.ListUsers()
		if err != nil {
			t.Errorf(err.Error())
		}
		if len(users) != 3 {
			t.Errorf("Three users expected")
		}
	})
	t.Run("Successfull ListUsers with nil users", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{
				listUsersRequest: &request.Request{},
				listUsersRequestOutput: &cognitoidentityprovider.ListUsersOutput{
					Users: nil,
				},
			},
		)
		users, err := cp.ListUsers()
		if err != nil {
			t.Errorf(err.Error())
		}
		if len(users) != 0 {
			t.Errorf("Zero users expected")
		}
	})
	t.Run("Fail ListUsers", func(t *testing.T) {
		cp := NewCognitoHandler(
			"client",
			"userpool",
			&mockedCognitoClient{
				listUsersRequest: &request.Request{Error: expectedError},
			},
		)
		_, err := cp.ListUsers()

		if err != expectedError {
			t.Errorf("Expected error")
		}
	})
}
