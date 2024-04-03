package conn

import (
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
)

func (h *HttpHandler) handleReplConf(conn net.Conn, _ internal.Request) {
	_, err := conn.Write([]byte(encoder.NewSimpleString("OK")))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}

func (h *HttpHandler) handlePsync(conn net.Conn, _ internal.Request) {
	_, err := conn.Write([]byte(encoder.NewSimpleString("FULLRESYNC <REPL_ID> 0")))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
