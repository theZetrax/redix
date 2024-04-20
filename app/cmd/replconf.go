package cmd

import (
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	SUBCMD_GETACK CMD_TYPE = "GETACK"
)

func handleReplConf(opts CMD_OPTS, args []any) []byte {
	log.Println("Handling REPLCONF command", args)

	sub_cmd := args[0].(string)
	if sub_cmd == string(SUBCMD_GETACK) {
		return resp.EncodeBulkString(
			"REPLCONF ACK 0",
		)
	}
	return resp.EncodeSimpleString("OK")
}
