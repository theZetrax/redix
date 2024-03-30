package internal

import (
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
	"github.com/codecrafters-io/redis-starter-go/app/internal/parser"
)

type HttpHandler struct{}

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
		default:
			_, err := conn.Write([]byte(encoder.NewError("ERR unknown command '" + cmd + "'")))
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
