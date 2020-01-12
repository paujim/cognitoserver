package main

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
)

type mockedCloudWatchClient struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
}

func (m *mockedCloudWatchClient) CreateLogGroup(*cloudwatchlogs.CreateLogGroupInput) (*cloudwatchlogs.CreateLogGroupOutput, error) {
	return &cloudwatchlogs.CreateLogGroupOutput{}, nil
}
func (m *mockedCloudWatchClient) CreateLogStream(*cloudwatchlogs.CreateLogStreamInput) (*cloudwatchlogs.CreateLogStreamOutput, error) {
	return &cloudwatchlogs.CreateLogStreamOutput{}, nil
}

func (m *mockedCloudWatchClient) PutLogEvents(*cloudwatchlogs.PutLogEventsInput) (*cloudwatchlogs.PutLogEventsOutput, error) {
	return &cloudwatchlogs.PutLogEventsOutput{}, nil
}

func TestEventFormat(t *testing.T) {
	event := createEvent("LOG", "Test Message")
	data := strings.SplitN(*event.Message, " ", 4)
	prefix, dateStr, message := data[0], data[1], data[3]

	if prefix != "[LOG]" {
		t.Errorf("Expected other prefix format")
	}
	_, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		t.Errorf("Expected date format")
	}

	if message != "Test Message" {
		t.Errorf("Expected date format")
	}
}

func TestLog(t *testing.T) {
	buf := &strings.Builder{}
	console := log.New(buf, "", 0)
	awsLog := NewAWSLogger(
		"LOG",
		"GName",
		"SName",
		&mockedCloudWatchClient{},
		console,
	)

	awsLog.Log("TEST")
	time.Sleep(time.Millisecond * (CollectionFrequency * 2))
	close(awsLog.dataCn)

	data := buf.String()
	fmt.Print(data)
	if data == "" {
		t.Errorf("Expected data")
	}

}
