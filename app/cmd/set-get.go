package cmd

import (
	"github.com/codecrafters-io/redis-starter-go/app/repository"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func handleSet(o CMD_OPTS, args []any) []byte {
	key := args[0].(string)
	value := args[1].(string)
	o.Store.Set(key, value, repository.SetOptions{
		HasTimeout: false,
		Timeout:    0,
	})

	return resp.EncodeSimpleString("OK")
}

func handleGet(o CMD_OPTS, args []any) []byte {
	key := args[0].(string)
	value, error := o.Store.Get(key)

	if error != nil {
		return resp.EncodeSimpleError(error.Error())
	}

	return resp.EncodeBulkString(value)
}
