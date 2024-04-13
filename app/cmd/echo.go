package cmd

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func handleEcho(_ CMD_OPTS, args []any) []byte {
	resp_raw := ""

	for idx, arg := range args {
		resp_raw = resp_raw + arg.(string)
		if idx != len(args)-1 {
			resp_raw = resp_raw + " "
		}
	}

	return resp.EncodeBulkString(resp_raw)
}
