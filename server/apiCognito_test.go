package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"testing"
)

// mocks
type mockedCognitoClient struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	request                      *request.Request
	initiateAuthOutput           *cognitoidentityprovider.InitiateAuthOutput
	respondToAuthChallengeOutput *cognitoidentityprovider.RespondToAuthChallengeOutput
}

func (m *mockedCognitoClient) InitiateAuthRequest(*cognitoidentityprovider.InitiateAuthInput) (*request.Request, *cognitoidentityprovider.InitiateAuthOutput) {
	return m.request, m.initiateAuthOutput
}

func (m *mockedCognitoClient) RespondToAuthChallengeRequest(*cognitoidentityprovider.RespondToAuthChallengeInput) (*request.Request, *cognitoidentityprovider.RespondToAuthChallengeOutput) {
	return m.request, m.respondToAuthChallengeOutput
}

func TestGetTokens(t *testing.T) {
	authResult := &cognitoidentityprovider.AuthenticationResultType{
		AccessToken:  aws.String("ACCESS_TOKEN"),
		RefreshToken: aws.String("REFRESH_TOKEN"),
	}
	logger := NewAPILogger("[TEST] ")

	t.Run("Missing parameters", func(t *testing.T) {
		cp := NewCognitoParam(
			"region",
			"client",
			"userpool",
			&mockedCognitoClient{},
			logger,
		)
		_, _, err := cp.GetTokens(nil, nil)
		if err != ErrorInvalidMissingCredentials {
			t.Errorf("Expected error when nil parameters")
		}
	})

	t.Run("Test successfull NEW_PASSWORD_REQUIRED challenge", func(t *testing.T) {
		cp := NewCognitoParam(
			"region",
			"client",
			"userpool",
			&mockedCognitoClient{
				request: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
					ChallengeName:        aws.String("NEW_PASSWORD_REQUIRED"),
					AuthenticationResult: authResult,
				},
				respondToAuthChallengeOutput: &cognitoidentityprovider.RespondToAuthChallengeOutput{
					AuthenticationResult: authResult,
				},
			},
			logger,
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
	t.Run("Test successfull challenge", func(t *testing.T) {
		cp := NewCognitoParam(
			"region",
			"client",
			"userpool",
			&mockedCognitoClient{
				request: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
					ChallengeName:        aws.String("OTHER"),
					AuthenticationResult: authResult,
				},
				respondToAuthChallengeOutput: &cognitoidentityprovider.RespondToAuthChallengeOutput{
					AuthenticationResult: authResult,
				},
			},
			logger,
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
	t.Run("Test error on challenge", func(t *testing.T) {

		expectedError := errors.New("Something went wrong")
		cp := NewCognitoParam(
			"region",
			"client",
			"userpool",
			&mockedCognitoClient{
				request:            &request.Request{Error: expectedError},
				initiateAuthOutput: nil,
			},
			logger,
		)
		_, _, err := cp.GetTokens(aws.String("username"), aws.String("password"))

		if err != expectedError {
			t.Errorf("Expected error")
		}

	})
}
