package cmd

import "github.com/codecrafters-io/redis-starter-go/app/resp"

func handleWait(opts CMD_OPTS, args []any) []byte {
	println("Handling Wait: ", opts.ConnectedNodeCount)
	return resp.EncodeInteger(opts.ConnectedNodeCount)
}
