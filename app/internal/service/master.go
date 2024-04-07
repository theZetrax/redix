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

type MainNodeOptions struct {
	IsMaster    bool
	ShouldClose bool
}
type MainNodeHandlerFn func(conn net.Conn, req internal.Request, opts MainNodeOptions)

// MainNode handles incoming requests
// and delegates them to the appropriate handler
type MainNode struct {
	StorageEngine *repository.StorageEngine
	Config        *internal.Config
	ConnPool      map[string]net.Conn // active connection pool
}

func (h *MainNode) Handle(
	conn net.Conn,
	buf *[]byte, // buffer
	read int, // read bytes length
	opts MainNodeOptions, // handler options
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
			continue READLOOP
		}

		// parse the request
		req, err := internal.ParseRequest((*buf)[:read])
		if err != nil {
			log.Println("Error parsing request: ", err.Error())
			return
		}
		log.Println(req.CMD.CMD, req.CMD.Args)

		// find the handler for the request
		var handler MainNodeHandlerFn = nil
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

		if handler == nil {
			return
		}

		// delegate the request to all connected replicas
		if opts.IsMaster && IsDelegateReq(req.CMD) && len(h.ConnPool) > 0 {
			for cid, conn := range h.ConnPool {
				_, err := conn.Write(*buf)
				if err != nil {
					conn.Close()
					delete(h.ConnPool, cid)
				}
			}

			log.Printf("delegated to active nodes[ACTIVE: %v]: %v %v", len(h.ConnPool), req.CMD.CMD, req.CMD.Args)
		}

		// handle the request
		handler(
			conn,
			req,
			opts,
		)
	}
}

// Close closes all active connections to the server
// and cleans up resources
func (h *MainNode) Close() {
	for _, conn := range h.ConnPool {
		conn.Close()
	}
}

// AddToConnPool adds a connection to the active connection pool
func (h *MainNode) AddToConnPool(conn net.Conn, uuid string) {
	h.ConnPool[uuid] = conn
}

func (h *MainNode) handleEcho(
	conn net.Conn, req internal.Request, _ MainNodeOptions,
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

func (h *MainNode) handlePing(
	conn net.Conn, req internal.Request, _ MainNodeOptions,
) {
	_, err := conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
