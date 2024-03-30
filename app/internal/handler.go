package internal

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
	"github.com/codecrafters-io/redis-starter-go/app/internal/parser"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
)

type HttpHandler struct {
	StorageEngine *repository.StorageEngine
}

func (h *HttpHandler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	var readErr error
	for readErr != io.EOF {
		buf := make([]byte, 1024)

		// request buffer length
		var rbLen int
		rbLen, readErr = conn.Read(buf)

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			log.Println("Error reading from connection: ", readErr.Error())
			os.Exit(1)
		}

		req := ParseRequest(buf[:rbLen])

		log.Println(req.CMD.CMD, req.CMD.Args)

		switch cmd := req.CMD.CMD; cmd {
		case parser.CMD_PING:
			_, err := conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				log.Println("Error writing to connection: ", err.Error())
				os.Exit(1)
			}
			break
		case parser.CMD_ECHO:
			h.handleEcho(conn, req)
			break
		case parser.CMD_SET:
			h.handleSet(conn, req)
			break
		case parser.CMD_GET:
			h.handleGet(conn, req)
			break
		default:
			log.Println("Unknown command: ", cmd)
			_, err := conn.Write([]byte(encoder.NewError(errors.New("ERR unknown command '" + cmd + "'"))))
			if err != nil {
				log.Println("Error writing to connection: ", err.Error())
				os.Exit(1)
			}
		}
	}
}

func (h *HttpHandler) handleEcho(conn net.Conn, req Request) {
	args_raw := req.CMD.Args
	args := encoder.ConvertSliceToStringArray(args_raw)
	resp := encoder.NewBulkString(strings.Join(args, " "))

	_, err := conn.Write([]byte(resp))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}

func (h *HttpHandler) handleSet(conn net.Conn, req Request) {
	args_raw := req.CMD.Args
	args := encoder.ConvertSliceToStringArray(args_raw)
	key := args[0]
	value := args[1:]

	if len(args) < 2 {
		_, err := conn.Write([]byte(encoder.NewError(errors.New("ERR wrong number of arguments for 'set' command"))))
		if err != nil {
			log.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
		return
	}

	h.StorageEngine.Set(key, strings.Join(value, " "))

	_, err := conn.Write([]byte(encoder.NewSimpleString("OK")))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
	}
}

func (h *HttpHandler) handleGet(conn net.Conn, req Request) {
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
