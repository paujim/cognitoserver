package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
)

//LOGFormat ... [Prefix] Date [Caller] Message
const LOGFormat = "[%s] %s [%s] %s"

//CollectionFrequency ... ms
const CollectionFrequency = 500

//AWSLogger ...
type AWSLogger struct {
	prefix     string
	groupName  *string
	streamName *string
	nextToken  *string
	err        error
	client     cloudwatchlogsiface.CloudWatchLogsAPI
	dataCn     chan cloudwatchlogs.InputLogEvent
	console    *log.Logger
}

//Log ...
func (l *AWSLogger) Log(v ...interface{}) {
	if l.err == nil {
		l.dataCn <- createEvent(l.prefix, fmt.Sprint(v...))
	}
}

//Logln ...
func (l *AWSLogger) Logln(v ...interface{}) {
	if l.err == nil {
		l.dataCn <- createEvent(l.prefix, fmt.Sprintln(v...))
	}
}

//Logf ..
func (l *AWSLogger) Logf(format string, v ...interface{}) {
	if l.err == nil {
		l.dataCn <- createEvent(l.prefix, fmt.Sprintf(format, v...))
	}
}

//NewAWSLogger ...
func NewAWSLogger(prefix, groupName, streamName string, client cloudwatchlogsiface.CloudWatchLogsAPI, console *log.Logger) *AWSLogger {

	awslog := &AWSLogger{
		prefix:     prefix,
		client:     client,
		groupName:  aws.String(groupName),
		streamName: aws.String(streamName),
		nextToken:  nil,
		err:        nil,
		dataCn:     make(chan cloudwatchlogs.InputLogEvent, 100),
		console:    console,
	}
	awslog.startCollectAndSendBatch()
	return awslog
}

func getCallerFunctionName() string {
	targetFrameIndex := 4
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame.Function
}

func createEvent(prefix, message string) cloudwatchlogs.InputLogEvent {
	caller := getCallerFunctionName()
	now := time.Now()
	timestamp := aws.Int64(int64(time.Nanosecond) * now.UnixNano() / int64(time.Millisecond))
	msg := aws.String(fmt.Sprintf(LOGFormat, prefix, now.Format(time.RFC3339), caller, message))
	event := cloudwatchlogs.InputLogEvent{
		Timestamp: timestamp,
		Message:   msg,
	}
	return event
}

func (l *AWSLogger) sendBatch(events []*cloudwatchlogs.InputLogEvent) ([]*cloudwatchlogs.InputLogEvent, int) {

	Ok := func(input *cloudwatchlogs.PutLogEventsInput, resp *cloudwatchlogs.PutLogEventsOutput, err error) ([]*cloudwatchlogs.InputLogEvent, int) {
		for _, event := range input.LogEvents {
			l.console.Println(*event)
		}
		l.nextToken = resp.NextSequenceToken
		l.err = err
		return []*cloudwatchlogs.InputLogEvent{}, 0
	}

	if len(events) == 0 || l.err != nil {
		return events, len(events)
	}

	input := &cloudwatchlogs.PutLogEventsInput{
		LogEvents:     events,
		SequenceToken: l.nextToken,
		LogGroupName:  l.groupName,
		LogStreamName: l.streamName,
	}
	resp, err := l.client.PutLogEvents(input)

	if err == nil {
		return Ok(input, resp, err)
	}

	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case cloudwatchlogs.ErrCodeResourceNotFoundException:
			l.client.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
				LogGroupName: l.groupName,
			})
			_, err = l.client.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
				LogGroupName:  l.groupName,
				LogStreamName: l.streamName,
			})
		case cloudwatchlogs.ErrCodeInvalidSequenceTokenException:
			res := strings.Split(aerr.Message(), " ")
			next := res[len(res)-1]
			if next == "null" {
				input.SequenceToken = nil
			} else {
				input.SetSequenceToken(next)
			}
			resp, err := l.client.PutLogEvents(input)
			if err == nil {
				return Ok(input, resp, err)
			}
		}
	}

	l.err = err
	return events, len(events)
}

/*
* The maximum batch size is 1,048,576 bytes, and this size is calculated
as the sum of all event messages in UTF-8, plus 26 bytes for each log
event.

* None of the log events in the batch can be more than 2 hours in the
future.

* None of the log events in the batch can be older than 14 days or older
than the retention period of the log group.

* The log events in the batch must be in chronological ordered by their
timestamp. The timestamp is the time the event occurred, expressed as
the number of milliseconds after Jan 1, 1970 00:00:00 UTC. (In AWS Tools
for PowerShell and the AWS SDK for .NET, the timestamp is specified in
.NET format: yyyy-mm-ddThh:mm:ss. For example, 2017-09-15T13:45:30.)

* A batch of log events in a single request cannot span more than 24 hours.
Otherwise, the operation fails.

* The maximum number of log events in a batch is 10,000.

* There is a quota of 5 requests per second per log stream. Additional
requests are throttled. This quota can't be changed.
*/
func (l *AWSLogger) startCollectAndSendBatch() {
	if l.err != nil {
		return
	}

	ticker := time.NewTicker(time.Millisecond * CollectionFrequency)
	go func() {
		batch := []*cloudwatchlogs.InputLogEvent{}
		size := 0
		for {
			select {
			// case <-done:
			// 	close(l.dataCn)
			case <-ticker.C:
				batch, size = l.sendBatch(batch)
			case event, more := <-l.dataCn:
				if !more {
					batch, size = l.sendBatch(batch)
					return
				}
				messageSize := len(*event.Message) + 26
				if size+messageSize >= 1048576 || len(batch) == 10000 {
					batch, size = l.sendBatch(batch)
				}
				batch = append(batch, &event)
				size += messageSize
			}
		}
	}()
}
