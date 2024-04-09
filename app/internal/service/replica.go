package service

import (
	"fmt"
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/decoder"
	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
	"github.com/codecrafters-io/redis-starter-go/app/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
)

// ReplicaNode handles the response from the master node
type ReplicaNode struct {
	StorageEngine *repository.StorageEngine
}

func (h *ReplicaNode) Handle(buf []byte) {
	req, err := internal.ParseRequest(buf)
	if err != nil {
		log.Println("Error parsing request: ", err.Error())
		return
	}

	fmt.Println("[REPLICA]", req.CMD.CMD, req.CMD.Args)
	logger.Log(logger.LOG_REQRES, fmt.Sprintf("[REPLICA] %s %v", req.CMD.CMD, req.CMD.Args))

	switch req.CMD.CMD {
	case decoder.CMD_SET:
		// handle set command
		h.handleSet(req)
	}
}

func (h *ReplicaNode) handleSet(req internal.Request) {
	args_raw := req.CMD.Args
	args := encoder.ConvertSliceToStringArray(args_raw)

	log.Println("Handling set command: ", args)
	key, values, opts, err := decoder.ParseSetCommand(args)
	if err != nil {
		log.Println("Error parsing set command: ", err.Error())
		return
	}

	h.StorageEngine.Set(key, values, opts)
	log.Println("Set command executed successfully")
}
