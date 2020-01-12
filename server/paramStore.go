package main

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

//ParamStore ...
type ParamStore struct {
	client ssmiface.SSMAPI
	logger ILogger
}

//NewParamStore ...
func NewParamStore(region string, logger ILogger) *ParamStore {
	mySession := session.Must(session.NewSession())
	store := &ParamStore{
		client: ssm.New(mySession, aws.NewConfig().WithRegion(region)),
		logger: logger,
	}
	return store
}

//Get ...
func (p *ParamStore) Get(key string) (string, error) {
	p.logger.Logf("Geting parameter [%v]\n", key)
	withDecryption := false
	input := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: &withDecryption,
	}
	param, err := p.client.GetParameter(input)
	if err != nil {
		p.logger.Logf("Fail: %v\n", err.Error())
		return "", err
	}
	if param.Parameter.Value == nil {
		p.logger.Logln("Fail: not found\n")
		return "", errors.New("not found")
	}
	return *param.Parameter.Value, nil
}
