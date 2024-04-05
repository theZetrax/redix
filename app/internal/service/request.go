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

type RequestHandlerOptions struct {
	IsMaster    bool
	ShouldClose bool
}
type ReqHandlerFunc func(conn net.Conn, req internal.Request, opts RequestHandlerOptions)

// ReqHandler handles incoming requests
// and delegates them to the appropriate handler
type ReqHandler struct {
	StorageEngine *repository.StorageEngine
	Config        *internal.Config
	ConnPool      map[string]net.Conn // active connection pool
}

func (h *ReqHandler) HandleRequest(
	conn net.Conn,
	buf *[]byte, // buffer
	read int, // read bytes length
	opts RequestHandlerOptions, // handler options
) {
	// by default should close the connection
	if opts.ShouldClose {
		defer conn.Close()
	}

	var readErr error
	init_loop := true // initial loop
READLOOP:
	for readErr != io.EOF {
		if !init_loop {
			read, readErr = conn.Read(*buf)
		}
		init_loop = false

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			log.Println("Error reading from connection: ", readErr.Error())
			return
		}

		// parse the request
		req := internal.ParseRequest((*buf)[:read])
		log.Println(req.CMD.CMD, req.CMD.Args)

		// find the handler for the request
		var handler ReqHandlerFunc = nil
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
				break READLOOP // break the loop
			}
		}

		// handle the request
		if handler == nil {
			return
		}

		if IsDelegateReq(req.CMD) {
			for _, conn := range h.ConnPool {
				conn.Write(*buf)
			}
		}

		handler(
			conn,
			req,
			opts,
		)
	}
}

// Close closes all active connections to the server
// and cleans up resources
func (h *ReqHandler) Close() {
	for _, conn := range h.ConnPool {
		conn.Close()
	}
}

// AddToConnPool adds a connection to the active connection pool
func (h *ReqHandler) AddToConnPool(conn net.Conn, uuid string) {
	h.ConnPool[uuid] = conn
}

func (h *ReqHandler) handleEcho(
	conn net.Conn, req internal.Request, _ RequestHandlerOptions,
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

func (h *ReqHandler) handlePing(
	conn net.Conn, req internal.Request, _ RequestHandlerOptions,
) {
	_, err := conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
