package cmd

import "github.com/codecrafters-io/redis-starter-go/app/resp"

func handleWait(opts CMD_OPTS, args []any) []byte {
	return resp.EncodeInteger(0)
}
