package list

import (
	"fmt"
	"log"
	"time"

	"github.com/jpignata/cwl/pkg/cloudwatchlogs"
)

type List struct {
	GroupName  string
	MaxResults int
}

func (l List) Run() {
	if len(l.GroupName) > 0 {
		l.listLogStreams()
	} else {
		l.listLogGroups()
	}
}

func (l List) listLogGroups() {
	client := cloudwatchlogs.NewClient()
	groupNames, err := client.GetLogGroupNames(l.MaxResults)

	if err != nil {
		log.Fatal(err)
	}

	for _, groupName := range groupNames {
		fmt.Println(groupName)
	}
}

func (l List) listLogStreams() {
	client := cloudwatchlogs.NewClient()
	streams, err := client.GetLogStreamNames(l.GroupName, l.MaxResults)

	if err != nil {
		log.Fatal(err)
	}

	for _, stream := range streams {
		lastEvent := time.Now().Sub(stream.LastEvent).Truncate(time.Second)

		fmt.Printf("%s (%s)\n", stream.Name, lastEvent)
	}
}
