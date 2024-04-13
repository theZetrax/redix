package cmd

import "github.com/codecrafters-io/redis-starter-go/app/resp"

func handleReplConf(opts CMD_OPTS, args []any) []byte {
	return resp.EncodeSimpleString("OK")
}
