package internal

import (
	"io"
	"log"
	"net"
	"os"
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
		// log the current request details
		log.Println("Request: ", req.CMD)

		_, err := conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			log.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	}
}
