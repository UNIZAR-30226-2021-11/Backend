package main

import (
	"Backend/pkg/events"
	"Backend/pkg/state"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

const (
	NUM_CLIENT = 1
)

func main() {

	var clients []Client

	for i := 0; i < NUM_CLIENT; i++ {
		c := Client{Id: uint32(40)}
		c.Start()
		clients = append(clients, c)
	}

	clients[0].CreateGame(1)

	//for i := 1; i < NUM_CLIENT; i++ {
	//	clients[i].JoinGame(1)
	//}

	time.Sleep(time.Second * 10)
	clients[0].PlayCard()
}

type Client struct {
	*websocket.Conn
	Id uint32
}

func (c *Client) Start() {
	c.Conn = newWsConn()

	c.WriteMessage(websocket.TextMessage, c.Id)

	// Receive messages
	go func() {
		defer func() {
			err := c.Close()
			if err != nil {
				log.Printf("%v", err)
			}
		}()
		for {
			_, message, _ := c.ReadMessage()
			log.Printf("Client %v:Message received: %s", c.Id, message)
		}
	}()
}

func (c *Client) JoinGame(game uint32) {
	event := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		EventType: 1,
	}
	_ = c.WriteJSON(event)
}

func (c *Client) CreateGame(game uint32) {
	event := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		EventType: 0,
	}
	_ = c.WriteJSON(event)
}

func (c *Client) PlayCard() {
	event := events.Event{
		GameID:    1,
		PlayerID:  c.Id,
		EventType: 3,
		Card:      state.CreateCard(state.SUIT1, 1),
	}
	_ = c.WriteJSON(event)
}

func newWsConn() *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: ":9000", Path: "/simulation"}

	// Establish connection
	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

	return c
}
