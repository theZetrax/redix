package cmd

import (
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func handleInfo(opts CMD_OPTS, _ []any) []byte {
	response := resp.EncodeNodeInfo(*opts.ReplicaInfo)

	log.Println("Response: ", string(response))

	return response
}
