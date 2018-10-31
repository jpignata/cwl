# cwl

`cwl` is a tiny CLI for CloudWatch Logs

## Usage

### list

`cwl list` returns all of your log groups. and `cwl list [group name]` will return the log streams within that group.

### get

`cwl get <group name>` returns log messages from a log group. You can optionally pass:

| long flag | short flag | description |
| --- | --- | --- |
| **--follow** | **-f** | Poll logs and continuously print new events |
| **--streams** | **-s** | Show logs from a specific stream (can be specified multiple times) |
| **--filter** | | Filter pattern to apply to the logs |
| **--start** | | Earliest time to return logs (e.g. -1h, 2018-01-01 09:36:00 EST) |
| **--end** | | Latest time to return logs (e.g. 2021-01-20 12:00:00 EST) |