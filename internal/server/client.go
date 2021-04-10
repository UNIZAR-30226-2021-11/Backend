package server

import (
	v1 "Backend/internal/server/v1"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const channelBufSize = 100

// Client struct holds client-specific variables.
type Client struct {
	id     uint32
	ws     *websocket.Conn
	ch     chan *[]byte
	doneCh chan bool
	sr 		*v1.SimulationRouter
}

// NewClient initializes a new Client struct with given websocket.
func NewClient(ws *websocket.Conn, sr *v1.SimulationRouter) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	ch := make(chan *[]byte, channelBufSize)
	doneCh := make(chan bool)

	return &Client{sr.idManager.NextPlayerId(), ws, ch, doneCh, sr}
}

// Conn returns client's websocket.Conn struct.
func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

// SendMessage sends game state to the client.
func (c *Client) SendMessage(bytes *[]byte) {
	select {
	case c.ch <- bytes:
	default:
		c.sr.monitor.AddDroppedMessage()
	}
}

// Done sends done message to the Client which closes the conection.
func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	defer func() {
		err := c.ws.Close()
		if err != nil {
			log.Println("Error:", err.Error())
		}
	}()

	log.Println("Listening write to client")
	for {
		select {

		case bytes := <-c.ch:
			before := time.Now()
			err := c.ws.WriteMessage(websocket.BinaryMessage, *bytes)
			after := time.Now()

			if err != nil {
				log.Println(err)
			} else {
				elapsed := after.Sub(before)
				c.sr.monitor.AddSendTime(elapsed)
			}

		case <-c.doneCh:
			c.doneCh <- true
			return
		}
	}
}

func (c *Client) listenRead() {
	defer func() {
		err := c.ws.Close()
		if err != nil {
			log.Println("Error:", err.Error())
		}
	}()

	log.Println("Listening read from client")
	for {
		select {

		case <-c.doneCh:
			c.doneCh <- true
			return

		default:
			c.readFromWebSocket()
		}
	}
}

func (c *Client) readFromWebSocket() {
	messageType, data, err := c.ws.ReadMessage()
	if err != nil {
		log.Println(err)

		c.doneCh <- true
		c.server.eventsDispatcher.FireUserLeft(&events.UserLeft{ClientID: c.id})
	} else if messageType != websocket.BinaryMessage {
		log.Println("Non binary message recived, ignoring")
	} else {
		c.unmarshalUserInput(data)
	}
}

