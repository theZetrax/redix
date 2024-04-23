package cmd

import (
	"fmt"
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
		if opts.ReplicaInfo.OffsetCount == -1 {
			opts.ReplicaInfo.OffsetCount = 0
		}
		offsetCount := fmt.Sprintf("%d", opts.ReplicaInfo.OffsetCount)

		return resp.EncodeArray(
			resp.EncodeBulkString(
				"REPLCONF",
			),
			resp.EncodeBulkString(
				"ACK",
			),
			resp.EncodeBulkString(
				offsetCount,
			),
		)
	}
	return resp.EncodeSimpleString("OK")
}
