package manager

import (
	"io"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/cmd"
	"github.com/codecrafters-io/redis-starter-go/app/logger"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte // broadcast message to all clients
	register   chan *Client
	unregister chan *Client
	store      *repository.Store
}

type Client struct {
	manager *ClientManager
	conn    net.Conn
	message chan []byte // Incoming requests from the clients.
	send    chan []byte // Outgoing responses to the clients.
}

func NewClientManager(store *repository.Store) *ClientManager {
	cm := &ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		store:      store,
	}

	return cm
}

func (cm *ClientManager) setup() {
	for {
		select {
		case client := <-cm.register:
			cm.clients[client] = true
		case client := <-cm.unregister:
			if _, ok := cm.clients[client]; ok {
				//TODO: close client connection
				delete(cm.clients, client)
				close(client.send)
				close(client.message)
			}
		case message := <-cm.broadcast:
			for client := range cm.clients {
				client.send <- message
			}
		}
	}
}

func NewClient(manager *ClientManager, conn net.Conn) *Client {
	c := &Client{
		conn:    conn,
		manager: manager,
		message: make(chan []byte),
	}

	return c
}

func (c *Client) Setup() {
	for {
		select {
		case message := <-c.message:
			if len(message) == 0 {
				continue
			}

			logger.LogResp("Received: ", message)
			handler, _ := resp.HandleResp(message)
			var response []byte

			switch handler.(type) {
			case *resp.Array:
				arr := handler.(*resp.Array)
				cmd := cmd.NewCMD(arr.Parsed, cmd.CMD_OPTS{
					Store: c.manager.store,
				})
				response = cmd.Process()
			default:
				response = handler.Process()
			}

			c.conn.Write(response)
		}
	}
}

func (c *Client) Read() {
	for {
		buf := make([]byte, 1024)
		read, err := c.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				c.manager.unregister <- c
				break
			}
		}

		if read == 0 {
			continue
		}

		message := buf[:read]
		c.message <- message
	}
}
