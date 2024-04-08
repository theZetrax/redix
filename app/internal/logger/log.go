package logger

import (
	"log"
	"os"
	"path"
)

const (
	LOG_REQRES = "REQ_RES"
	LOG_ERR    = "ERR"
)

const log_path = "tmp"

// Log writes a string to a file
// for logging purposes
func Log(log_type string, str string) {
	if _, err := os.Stat(log_path); os.IsNotExist(err) {
		err = os.MkdirAll(log_path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.OpenFile(path.Join(log_path, "LOGS.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	logger := log.New(f, log_type, log.LstdFlags)
	logger.Println(str)
}
