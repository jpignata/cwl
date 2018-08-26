package cloudwatchlogs

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awscwl "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type Client struct {
	client *awscwl.CloudWatchLogs
}

type GetLogsInput struct {
	EndTime     time.Time
	Filter      string
	GroupName   string
	StartTime   time.Time
	StreamNames []string
}

type Line struct {
	EventId    string
	GroupName  string
	StreamName string
	Message    string
	Timestamp  time.Time
}

type Lines []Line

type Stream struct {
	Name      string
	LastEvent time.Time
}

type Streams []Stream

func NewClient() Client {
	sess := session.Must(session.NewSession())

	return Client{
		client: awscwl.New(sess),
	}
}

func (c Client) GetLogs(i GetLogsInput) ([]Line, error) {
	lines := make(Lines, 0)
	input := &awscwl.FilterLogEventsInput{
		LogGroupName: aws.String(i.GroupName),
		Interleaved:  aws.Bool(true),
	}

	if !i.StartTime.IsZero() {
		input.SetStartTime(i.StartTime.UTC().UnixNano() / int64(time.Millisecond))
	}

	if !i.EndTime.IsZero() {
		input.SetEndTime(i.EndTime.UTC().UnixNano() / int64(time.Millisecond))
	}

	if len(i.Filter) > 0 {
		input.SetFilterPattern(i.Filter)
	}

	if len(i.StreamNames) > 0 {
		input.SetLogStreamNames(aws.StringSlice(i.StreamNames))
	}

	err := c.client.FilterLogEventsPages(
		input,
		func(resp *awscwl.FilterLogEventsOutput, lastPage bool) bool {
			for _, event := range resp.Events {
				lines = append(lines,
					Line{
						EventId:    aws.StringValue(event.EventId),
						Message:    aws.StringValue(event.Message),
						GroupName:  i.GroupName,
						StreamName: aws.StringValue(event.LogStreamName),
						Timestamp:  time.Unix(0, aws.Int64Value(event.Timestamp)*int64(time.Millisecond)),
					},
				)
			}

			return true
		},
	)

	return lines, err
}

func (c Client) GetLogGroupNames(maxResults int) ([]string, error) {
	results := 0
	groupNames := make([]string, 0)

	err := c.client.DescribeLogGroupsPages(
		&awscwl.DescribeLogGroupsInput{},
		func(resp *awscwl.DescribeLogGroupsOutput, lastPage bool) bool {
			for _, logGroup := range resp.LogGroups {
				groupNames = append(groupNames, aws.StringValue(logGroup.LogGroupName))
				results += 1
			}

			return results < maxResults
		},
	)

	return groupNames, err
}

func (c Client) GetLogStreamNames(logGroup string, maxResults int) (Streams, error) {
	results := 0
	streams := make(Streams, 0)
	input := &awscwl.DescribeLogStreamsInput{
		Descending:   aws.Bool(true),
		LogGroupName: aws.String(logGroup),
		OrderBy:      aws.String("LastEventTime"),
	}

	err := c.client.DescribeLogStreamsPages(
		input,
		func(resp *awscwl.DescribeLogStreamsOutput, lastPage bool) bool {
			for _, logStream := range resp.LogStreams {
				streams = append(streams,
					Stream{
						Name:      aws.StringValue(logStream.LogStreamName),
						LastEvent: time.Unix(0, aws.Int64Value(logStream.LastEventTimestamp)*int64(time.Millisecond)),
					},
				)

				results += 1
			}

			return results < maxResults
		},
	)

	return streams, err
}
