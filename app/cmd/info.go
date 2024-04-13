package cmd

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func handleInfo(opts CMD_OPTS, _ []any) []byte {
	return resp.EncodeNodeInfo(*opts.ReplicaInfo)
}
