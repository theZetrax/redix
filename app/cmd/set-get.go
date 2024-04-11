package cmd

import (
	"errors"
	"log"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/repository"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	SUB_PX = "PX"
)

var SUB_COMMANDS = map[CMD_TYPE][]string{
	CMD_SET: {SUB_PX},
}

func handleSet(o CMD_OPTS, _args []any) []byte {
	args := reflect.ValueOf(_args).Interface().([]string)
	key, value, opts, err := parseSetCmd(args)

	if err != nil {
		return resp.EncodeSimpleError(err.Error())
	}

	o.Store.Set(key, value, opts)
	return resp.EncodeSimpleString("OK")
}
func parseSetCmd(args []string) (string, string, repository.SetOptions, error) {
	if len(args) < 2 {
		return "", "", repository.SetOptions{}, errors.New("ERR wrong number of arguments for 'set' command")
	}

	key := args[0]
	values := args[1:]
	sub_cmd_map := make(map[string]any) // map to store sub commands
	cmd_start_index := -1

	for _, sub_cmd := range SUB_COMMANDS[CMD_SET] {
		for idx, value := range values {
			if strings.ToUpper(value) == sub_cmd {
				if idx+1 >= len(values) {
					return "", "", repository.SetOptions{}, errors.New("ERR wrong number of arguments for 'set' command")
				}

				cmd_start_index = idx
				sub_cmd_map[sub_cmd] = values[idx+1]
			}
		}

	}

	opts := repository.SetOptions{
		HasTimeout: false,
	}

	// handle sub commands
	if cmd_start_index != -1 {
		values = slices.Clone(values[:cmd_start_index])
		for k, v := range sub_cmd_map {
			switch k {
			case SUB_PX:
				// handle px sub command
				log.Print("Handling PX command: ", v)
				timeout, err := strconv.Atoi(reflect.ValueOf(v).String())
				if err != nil {
					return "", "", repository.SetOptions{}, errors.New("ERR wrong number of arguments for 'set' command")
				}
				opts.HasTimeout = true
				opts.Timeout = timeout
			}
		}
	}

	return key, strings.Join(values, " "), opts, nil

}

func handleGet(o CMD_OPTS, args []any) []byte {
	key := args[0].(string)
	value, error := o.Store.Get(key)

	if error != nil {
		return resp.EncodeSimpleError(error.Error())
	}

	return resp.EncodeBulkString(value)
}
