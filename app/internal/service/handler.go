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
	StorageEngine     *repository.StorageEngine
	Config            *internal.Config
	ActiveConnections map[string]net.Conn
	ShouldClose       bool
}

func (h *HttpHandler) HandleConnection(conn net.Conn) {
	// by default should close the connection
	h.ShouldClose = true

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

		// parse the request
		req := internal.ParseRequest(buf[:rbLen])
		log.Println(req.CMD.CMD, req.CMD.Args)

		// find the handler for the request
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
			err_resp := encoder.NewError(
				errors.New("ERR unknown command '" + cmd + "'"),
			)

			// write the error response
			_, err := conn.Write([]byte(err_resp))
			if err != nil {
				log.Println("Error writing to connection: ", err.Error())
				os.Exit(1)
			}
		}

		// handle the request
		if handler != nil {
			handler(
				conn,
				req,
				HandlerOptions{
					IsMaster: h.Config.IsMaster,
				},
			)
		}
	}

	// close the connection if the flag is set
	if h.ShouldClose {
		conn.Close()
	}
}

// Close closes all active connections to the server
// and cleans up resources
func (h *HttpHandler) Close() {
	for _, conn := range h.ActiveConnections {
		conn.Close()
	}
}

// AddToConnPool adds a connection to the active connection pool
func (h *HttpHandler) AddToConnPool(conn net.Conn, uuid string) {
	h.ShouldClose = false // do not close the connection
	h.ActiveConnections[uuid] = conn
}

func (h *HttpHandler) handleEcho(
	conn net.Conn, req internal.Request, _ HandlerOptions,
) {
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
