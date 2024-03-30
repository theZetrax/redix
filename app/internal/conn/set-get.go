package conn

import (
	"errors"
	"log"
	"net"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
	"github.com/codecrafters-io/redis-starter-go/app/internal/parser"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
)

func (h *HttpHandler) handleSet(conn net.Conn, req internal.Request) {
	args_raw := req.CMD.Args
	args := encoder.ConvertSliceToStringArray(args_raw)

	if len(args) < 2 {
		_, err := conn.Write([]byte(encoder.NewError(errors.New("ERR wrong number of arguments for 'set' command"))))
		if err != nil {
			log.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
		return
	}

	key := args[0]
	values := args[1:]
	sub_cmd_map := make(map[string]any) // map to store sub commands
	cmd_start_index := -1

	for _, sub_cmd := range parser.SUB_COMMANDS[parser.CMD_SET] {
		for idx, value := range values {
			if strings.ToUpper(value) == sub_cmd {
				if idx+1 >= len(values) {
					_, err := conn.Write([]byte(encoder.NewError(errors.New("ERR wrong number of arguments for 'set' command"))))
					if err != nil {
						log.Println("Error writing to connection: ", err.Error())
						os.Exit(1)
					}
					return // exit
				}

				cmd_start_index = idx
				sub_cmd_map[sub_cmd] = values[idx+1]
			}
		}

	}

	values = slices.Clone(values[:cmd_start_index])
	opts := repository.SetOptions{
		HasTimeout: false,
	}

	// handle sub commands
	if cmd_start_index != -1 {
		for k, v := range sub_cmd_map {
			switch k {
			case parser.SUB_PX:
				// handle px sub command
				log.Print("Handling PX command: ", v)
				timeout, err := strconv.Atoi(reflect.ValueOf(v).String())
				if err != nil {
					_, err = conn.Write([]byte(encoder.NewError(errors.New("ERR wrong number of arguments for 'set' command"))))
					if err != nil {
						log.Println("Error writing to connection: ", err.Error())
						os.Exit(1)
					}
					return
				}
				opts.HasTimeout = true
				opts.Timeout = timeout
				break
			}
		}
	}

	h.StorageEngine.Set(key, strings.Join(values, " "), opts)

	_, err := conn.Write([]byte(encoder.NewSimpleString("OK")))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
	}
}

func (h *HttpHandler) handleGet(conn net.Conn, req internal.Request) {
	args_raw := req.CMD.Args
	args := encoder.ConvertSliceToStringArray(args_raw)

	if len(args) != 1 {
		_, err := conn.Write([]byte(encoder.NewError(errors.New("ERR wrong number of arguments for 'get' command"))))
		if err != nil {
			log.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
		return
	}

	key := args[0]

	value, err := h.StorageEngine.Get(key)

	if err != nil {
		_, err := conn.Write([]byte(encoder.NewError(err)))
		if err != nil {
			log.Println("Error writing to connection: ", err.Error())
		}
		return
	}

	_, err = conn.Write([]byte(encoder.NewBulkString(value)))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
