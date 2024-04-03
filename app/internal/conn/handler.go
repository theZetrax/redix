package conn

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

		switch cmd := req.CMD.CMD; cmd {
		case decoder.CMD_PING:
			_, err := conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				log.Println("Error writing to connection: ", err.Error())
				os.Exit(1)
			}
			break
		case decoder.CMD_ECHO:
			h.handleEcho(conn, req)
			break
		case decoder.CMD_SET:
			h.handleSet(conn, req)
			break
		case decoder.CMD_GET:
			h.handleGet(conn, req)
			break
		case decoder.CMD_INFO:
			h.handleInfo(conn, req)
			break
		case decoder.CMD_REPLCONF:
			h.handleReplConf(conn, req)
			break
		case decoder.CMD_PSYNC:
			h.handlePsync(conn, req)
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

func (h *HttpHandler) handleEcho(conn net.Conn, req internal.Request) {
	args_raw := req.CMD.Args
	args := encoder.ConvertSliceToStringArray(args_raw)
	resp := encoder.NewBulkString(strings.Join(args, " "))

	_, err := conn.Write([]byte(resp))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
