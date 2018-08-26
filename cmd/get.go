package cmd

import (
	"github.com/jpignata/cwl/pkg/get"
	"github.com/spf13/cobra"
)

var (
	flagFilter    string
	flagEndTime   string
	flagStartTime string
	flagFollow    bool
	flagStreams   []string
)

var getCmd = &cobra.Command{
	Use:   "get <group name>",
	Short: "Show logs from a log group",
	Long: `Show logs from a log group

Return either a specific segment of logs or tail logs in real-time
using the --follow option.

Follow will continue to run and return logs until interrupted by Control-C. If
--follow is passed --end cannot be specified.

Logs can be returned for specific stream by passing the stream name via the
--stream flag. It can be passed multiple times for multiple log streams.

A specific window of logs can be requested by passing --start and --end options
with a time expression. The time expression can be either a duration or a
timestamp:

  - Duration (e.g. -1h [one hour ago], -1h10m30s [one hour, ten minutes, and
    thirty seconds ago], 2h [two hours from now])
  - Timestamp with optional timezone in the format of YYYY-MM-DD HH:MM:SS [TZ];
    timezone will default to UTC if omitted (e.g. 2017-12-22 15:10:03 EST)

You can filter logs for specific term by passing a filter expression via the
--filter flag. Pass a single term to search for that term, pass multiple terms
to search for log messages that include all terms.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		get := &get.Get{
			LogGroupName: args[0],
			Filter:       flagFilter,
			Follow:       flagFollow,
		}

		get.AddStreams(flagStreams)
		get.AddStartTime(flagStartTime)
		get.AddEndTime(flagEndTime)

		get.Run()
	},
}

func init() {
	getCmd.Flags().BoolVarP(&flagFollow, "follow", "f", false, "Poll logs and continuously print new events")
	getCmd.Flags().StringVar(&flagFilter, "filter", "", "Filter pattern to apply")
	getCmd.Flags().StringVar(&flagStartTime, "start", "", "Earliest time to return logs (e.g. -1h, 2018-01-01 09:36:00 EST")
	getCmd.Flags().StringVar(&flagEndTime, "end", "", "Latest time to return logs (e.g. 3y, 2021-01-20 12:00:00 EST")
	getCmd.Flags().StringSliceVarP(&flagStreams, "streams", "s", []string{}, "Show logs from specific stream (can be specified multiple times)")

	rootCmd.AddCommand(getCmd)
}
