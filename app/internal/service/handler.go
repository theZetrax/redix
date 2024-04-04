package service

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/decoder"
	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
)

type HandlerOptions struct {
	IsMaster bool
}
type HandlerFunc func(conn net.Conn, req internal.Request, opts HandlerOptions)

type HttpHandler struct {
	StorageEngine *repository.StorageEngine
	Config        *internal.Config
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
			return
		}

		req := internal.ParseRequest(buf[:rbLen])

		log.Println(req.CMD.CMD, req.CMD.Args)

		var handler HandlerFunc = nil

		switch cmd := req.CMD.CMD; cmd {
		case decoder.CMD_PING:
			handler = h.handlePing
		case decoder.CMD_ECHO:
			handler = h.handleEcho
		case decoder.CMD_SET:
			handler = h.handleSet
		case decoder.CMD_GET:
			handler = h.handleGet
		case decoder.CMD_INFO:
			handler = h.handleInfo
		case decoder.CMD_REPLCONF:
			handler = h.handleReplConf
		case decoder.CMD_PSYNC:
			handler = h.handlePsync
		default:
			log.Println("Unknown command: ", cmd)
			_, err := conn.Write([]byte(encoder.NewError(errors.New("ERR unknown command '" + cmd + "'"))))
			if err != nil {
				log.Println("Error writing to connection: ", err.Error())
				os.Exit(1)
			}
		}

		if handler != nil {
			// handle the request
			handler(
				conn,
				req,
				HandlerOptions{
					IsMaster: h.Config.IsMaster,
				},
			)
		}
	}
}

func (h *HttpHandler) handleEcho(conn net.Conn, req internal.Request, _ HandlerOptions) {
	args_raw := req.CMD.Args
	args := encoder.ConvertSliceToStringArray(args_raw)
	resp := encoder.NewBulkString(strings.Join(args, " "))

	_, err := conn.Write([]byte(resp))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}

func (h *HttpHandler) handlePing(conn net.Conn, req internal.Request, _ HandlerOptions) {
	_, err := conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
