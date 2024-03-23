package internal

import (
	"log"
	"net"
	"os"
)

type HttpHandler struct{}

func (h *HttpHandler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		conn.Write([]byte("-ERR internal error\r\n"))
	}

	_, err = conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
