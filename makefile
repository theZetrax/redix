BIN_NAME=server
BIN_DIR=bin
BIN_PATH=${BIN_DIR}/${BIN_NAME}

build:
	go build -o ${BIN_PATH} -v app/server.go

run:
	go run app/server.go

run-replica:
	go run app/server.go --port 6381 --replicaof localhost 6379

build-run:
	go build -o ${BIN_PATH} -v app/server.go
	./${BIN_PATH}

clean:
	go clean
	rm -f ${BIN_PATH}
