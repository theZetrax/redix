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
	len, err := conn.Read(buf)
	if err != nil {
		conn.Write([]byte("-ERR internal error\r\n"))
	}

	req := ParseRequest(buf[:len])
	// log incoming request details
	log.Println(req.Method, req.Url, req.HttpVersion)

	_, err = conn.Write([]byte("+PONG" + CLRF))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
