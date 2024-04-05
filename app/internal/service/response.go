package service

import (
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/decoder"
	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
)

// ResponseHandler handles the response from the master node
type ResponseHandler struct {
	StorageEngine *repository.StorageEngine
}

func (h *ResponseHandler) HandleResponse(buf []byte) {
	req := internal.ParseRequest(buf)

	log.Println("[FROM MASTER]", req.CMD.CMD, req.CMD.Args)

	switch req.CMD.CMD {
	case decoder.CMD_SET:
		// handle set command
		h.handleSet(req)
		break
	}
}

func (h *ResponseHandler) handleSet(req internal.Request) {
	args_raw := req.CMD.Args
	args := encoder.ConvertSliceToStringArray(args_raw)

	key, values, opts, err := decoder.ParseSetCommand(args)
	if err != nil {
		log.Println("Error parsing set command: ", err.Error())
		return
	}

	h.StorageEngine.Set(key, values, opts)
	log.Println("Set command executed successfully")
}
