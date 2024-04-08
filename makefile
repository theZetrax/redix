BIN_NAME=server
BIN_DIR=bin
BIN_PATH=${BIN_DIR}/${BIN_NAME}

run:
	go run app/server.go

replica:
	go run app/server.go --port 6381 --replicaof localhost 6379

build:
	go build -o ${BIN_PATH} app/server.go

clean:
	rm -rf ${BIN_DIR}
