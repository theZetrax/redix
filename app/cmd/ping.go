package cmd

import "github.com/codecrafters-io/redis-starter-go/app/resp"

func handlePing(_ CMD_OPTS, _ []any) []byte {
	return resp.EncodeSimpleString("PONG")
}
