package service

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/decoder"
	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
)

func (h *ReqHandler) handleSet(conn net.Conn, req internal.Request, _ RequestHandlerOptions) {
	args_raw := req.CMD.Args
	args := encoder.ConvertSliceToStringArray(args_raw)

	key, values, opts, err := decoder.ParseSetCommand(args)
	if err != nil {
		_, err := conn.Write([]byte(encoder.NewError(err)))
		if err != nil {
			log.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	}

	h.StorageEngine.Set(key, values, opts)

	_, err = conn.Write([]byte(encoder.NewSimpleString("OK")))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
	}
}

func (h *ReqHandler) handleGet(conn net.Conn, req internal.Request, _ RequestHandlerOptions) {
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

	// check if expired
	if h.StorageEngine.ExpiredTimeout(key) {
		_, err := conn.Write([]byte(encoder.NewNil()))
		if err != nil {
			log.Println("Error writing to connection: ", err.Error())
		}
		return
	}

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
