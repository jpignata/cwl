package get

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/jpignata/cwl/pkg/cloudwatchlogs"
	"github.com/mgutz/ansi"
)

const (
	timeOutputFormat = "Jan 02 15:04:05"

	timeFormat         = "2006-01-02 15:04:05"
	timeFormatWithZone = "2006-01-02 15:04:05 MST"
	eventCacheSize     = 10000
)

type Empty struct{}

type Get struct {
	LogGroupName    string
	EndTime         time.Time
	Filter          string
	Follow          bool
	LogStreamColors map[string]string
	LogStreamNames  []string
	StartTime       time.Time
	EventCache      *lru.Cache
}

func (g *Get) AddStartTime(rawStartTime string) {
	if rawStartTime != "" {
		g.StartTime = g.parseTime(rawStartTime)
	}
}

func (g *Get) AddEndTime(rawEndTime string) {
	if rawEndTime != "" {
		g.EndTime = g.parseTime(rawEndTime)
	}
}

func (g *Get) AddStreams(streamNames []string) {
	for _, streamName := range streamNames {
		g.LogStreamNames = append(g.LogStreamNames, streamName)
	}
}

func (g Get) Validate() {
	if g.Follow && !g.EndTime.IsZero() {
		log.Fatal("--end-time cannot be specified if following")
	}
}

func (g *Get) getStreamColor(logStreamName string) string {
	if g.LogStreamColors == nil {
		g.LogStreamColors = make(map[string]string)
	}

	if g.LogStreamColors[logStreamName] == "" {
		g.LogStreamColors[logStreamName] = strconv.Itoa(rand.Intn(256))
	}

	return g.LogStreamColors[logStreamName]
}

func (g *Get) SeenEvent(eventId string) bool {
	if g.EventCache == nil {
		g.EventCache, _ = lru.New(eventCacheSize)
	}

	if !g.EventCache.Contains(eventId) {
		g.EventCache.Add(eventId, Empty{})
		return false
	} else {
		return true
	}
}

func (g Get) parseTime(rawTime string) time.Time {
	var t time.Time

	if duration, err := time.ParseDuration(strings.ToLower(rawTime)); err == nil {
		return time.Now().Add(duration)
	}

	if t, err := time.Parse(timeFormat, rawTime); err == nil {
		return t
	}

	if t, err := time.Parse(timeFormatWithZone, rawTime); err == nil {
		return t
	}

	log.Fatalf("Could not parse %s", rawTime)

	return t
}

func (g *Get) Run() {
	rand.Seed(time.Now().UTC().UnixNano())

	if g.Follow {
		g.followLogs()
	} else {
		g.getLogs()
	}
}

func (g Get) followLogs() {
	ticker := time.NewTicker(time.Second)

	if g.StartTime.IsZero() {
		g.StartTime = time.Now()
	}

	for {
		g.getLogs()

		if newStartTime := time.Now().Add(-10 * time.Second); newStartTime.After(g.StartTime) {
			g.StartTime = newStartTime
		}

		<-ticker.C
	}
}

func (g *Get) getLogs() {
	client := cloudwatchlogs.NewClient()
	input := cloudwatchlogs.GetLogsInput{
		StreamNames: g.LogStreamNames,
		GroupName:   g.LogGroupName,
		Filter:      g.Filter,
		StartTime:   g.StartTime,
		EndTime:     g.EndTime,
	}

	lines, err := client.GetLogs(input)

	if err != nil {
		log.Println(err)
	}

	for _, logLine := range lines {
		if !g.SeenEvent(logLine.EventId) {
			streamColor := g.getStreamColor(logLine.StreamName)

			fmt.Printf("%s %s %s%s%s %s\n",
				logLine.Timestamp.Format(timeOutputFormat),
				logLine.GroupName,
				ansi.ColorCode(streamColor),
				logLine.StreamName,
				ansi.ColorCode("reset"),
				logLine.Message,
			)
		}
	}
}
