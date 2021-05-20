package services

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/paujim/cognitoserver/server/pkg/entities"
	log "github.com/sirupsen/logrus"
)

type parameterStore struct {
	ssmAPI ssmiface.SSMAPI
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func NewParameterStore(api ssmiface.SSMAPI) entities.ParameterStorer {
	return &parameterStore{
		ssmAPI: api,
	}
}

func (p *parameterStore) Get(key string) (string, error) {
	log.Infof("Geting parameter [%v]\n", key)
	withDecryption := false
	input := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: &withDecryption,
	}
	param, err := p.ssmAPI.GetParameter(input)
	if err != nil {
		log.Errorf("Fail: %v\n", err.Error())
		return "", err
	}
	if param.Parameter.Value == nil {
		log.Errorf("Fail: not found\n")
		return "", errors.New("not found")
	}
	return *param.Parameter.Value, nil
}
