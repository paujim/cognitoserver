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
	initiateAuthRequest           *request.Request
	initiateAuthOutput            *cognitoidentityprovider.InitiateAuthOutput
	respondToAuthChallengeRequest *request.Request
	respondToAuthChallengeOutput  *cognitoidentityprovider.RespondToAuthChallengeOutput
}

func (m *mockedCognitoClient) InitiateAuthRequest(*cognitoidentityprovider.InitiateAuthInput) (*request.Request, *cognitoidentityprovider.InitiateAuthOutput) {
	return m.initiateAuthRequest, m.initiateAuthOutput
}

func (m *mockedCognitoClient) RespondToAuthChallengeRequest(*cognitoidentityprovider.RespondToAuthChallengeInput) (*request.Request, *cognitoidentityprovider.RespondToAuthChallengeOutput) {
	return m.respondToAuthChallengeRequest, m.respondToAuthChallengeOutput
}

func TestGetTokens(t *testing.T) {
	authResult := &cognitoidentityprovider.AuthenticationResultType{
		AccessToken:  aws.String("ACCESS_TOKEN"),
		RefreshToken: aws.String("REFRESH_TOKEN"),
	}
	expectedError := errors.New("Something went wrong")
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
		if err != ErrorInvalidInputParameters {
			t.Errorf("Expected error when nil parameters")
		}
	})
	t.Run("Test successfull NEW_PASSWORD_REQUIRED challenge", func(t *testing.T) {
		cp := NewCognitoParam(
			"region",
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
	t.Run("Test error after NEW_PASSWORD_REQUIRED challenge", func(t *testing.T) {
		cp := NewCognitoParam(
			"region",
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
					ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
				},
				respondToAuthChallengeRequest: &request.Request{Error: expectedError},
			},
			logger,
		)
		_, _, err := cp.GetTokens(aws.String("username"), aws.String("password"))
		if err != expectedError {
			t.Errorf("Ëxpected error")
		}
	})
	t.Run("Test OTHER challenge", func(t *testing.T) {
		cp := NewCognitoParam(
			"region",
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
					ChallengeName:        aws.String("OTHER"),
					AuthenticationResult: authResult,
				},
			},
			logger,
		)
		_, _, err := cp.GetTokens(aws.String("username"), aws.String("password"))
		if err == nil {
			t.Errorf("Error expected")
		}
	})
	t.Run("Test successfull challenge", func(t *testing.T) {
		cp := NewCognitoParam(
			"region",
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
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
	t.Run("Test error on request", func(t *testing.T) {
		cp := NewCognitoParam(
			"region",
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{Error: expectedError},
				initiateAuthOutput:  nil,
			},
			logger,
		)
		_, _, err := cp.GetTokens(aws.String("username"), aws.String("password"))

		if err != expectedError {
			t.Errorf("Expected error")
		}
	})
}

func TestRefreshAccessToken(t *testing.T) {
	authResult := &cognitoidentityprovider.AuthenticationResultType{
		AccessToken:  aws.String("ACCESS_TOKEN"),
		RefreshToken: aws.String("REFRESH_TOKEN"),
	}
	expectedError := errors.New("Something went wrong")
	logger := NewAPILogger("[TEST] ")
	t.Run("Missing parameters", func(t *testing.T) {
		cp := NewCognitoParam(
			"region",
			"client",
			"userpool",
			&mockedCognitoClient{},
			logger,
		)
		_, _, err := cp.RefreshAccessToken(nil)
		if err != ErrorInvalidInputParameters {
			t.Errorf("Expected error when nil parameters")
		}
	})

	t.Run("Test successfull refresh AccessToken", func(t *testing.T) {
		cp := NewCognitoParam(
			"region",
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
					AuthenticationResult: authResult,
				},
			},
			logger,
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
	t.Run("Test OTHER challenge", func(t *testing.T) {
		cp := NewCognitoParam(
			"region",
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{},
				initiateAuthOutput: &cognitoidentityprovider.InitiateAuthOutput{
					ChallengeName: aws.String("OTHER"),
				},
			},
			logger,
		)
		_, _, err := cp.RefreshAccessToken(aws.String("refresh_token"))
		if err == nil {
			t.Errorf("Error expected")
		}
	})
	t.Run("Test error on request", func(t *testing.T) {
		cp := NewCognitoParam(
			"region",
			"client",
			"userpool",
			&mockedCognitoClient{
				initiateAuthRequest: &request.Request{Error: expectedError},
			},
			logger,
		)
		_, _, err := cp.RefreshAccessToken(aws.String("refresh_token"))

		if err != expectedError {
			t.Errorf("Expected error")
		}
	})
}
