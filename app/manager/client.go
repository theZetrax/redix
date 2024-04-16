package manager

import (
	"io"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/cmd"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// ClientManager manages the clients connected to the server.
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte // broadcast message to all clients
	register   chan *Client
	unregister chan *Client
	store      *repository.Store
	node_info  *resp.NodeInfo
}

type Client struct {
	manager *ClientManager
	conn    net.Conn
	message chan []byte // Incoming requests from the clients.
	send    chan []byte // Outgoing responses to the clients.
}

func NewClientManager(store *repository.Store, node_info *resp.NodeInfo) *ClientManager {
	cm := &ClientManager{
		clients:    make(map[*Client]bool, 0),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		store:      store,
		node_info:  node_info,
	}

	return cm
}

func (cm *ClientManager) setup() {
	for {
		select {
		case client := <-cm.register:
			log.Println("Replica registered to master")
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
		send:    make(chan []byte),
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

			handler, _ := resp.HandleResp(message)
			var response []byte

			switch handler.(type) {
			case *resp.Array:
				arr := handler.(*resp.Array)
				cmd_handler := cmd.NewCMD(arr.Parsed, cmd.CMD_OPTS{
					Store:       c.manager.store,
					ReplicaInfo: c.manager.node_info,
				})

				switch {
				case cmd_handler.Name == cmd.CMD_PSYNC: // Register the replica to the master node.
					cmd_handler.Process(&c.conn, func() {
						c.manager.register <- c
					})
				case cmd_handler.Name == cmd.CMD_SET:
					cmd_handler.Process(&c.conn, nil)
					c.manager.broadcast <- message
				default:
					cmd_handler.Process(&c.conn, nil)
				}

				log.Println("Connections: ", c.manager.clients)
			default:
				response = handler.Process()
				c.conn.Write(response)
			}
		case message := <-c.send:
			if len(message) == 0 {
				continue
			}

			c.conn.Write(message)
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
