package cmd

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/utl"
)

func handlePsync(_ CMD_OPTS, args []any) [][]byte {
	id := "1234567890"
	offset := "0"
	response := make([][]byte, 0)

	f_resync_resp := resp.EncodeSimpleString(fmt.Sprintf("FULLRESYNC %s %s", id, offset))
	response = append(response, f_resync_resp)

	rdb_hex := []byte("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
	rdb_enc, err := utl.DecodeHexToBinary(rdb_hex)
	if err != nil {
		panic("Failed to decode hex to binary")
	}

	response = append(response, resp.EncodeFileContent(rdb_enc))

	return response
}
