package main

import (
	"Backend/pkg/events"
	"Backend/pkg/simulation"
	"Backend/pkg/state"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

const (
	NUM_CLIENT = 4
)

func main() {

	var clients []Client

	for i := 0; i < NUM_CLIENT; i++ {
		c := Client{Id: uint32(40 + i)}
		c.Start()
		clients = append(clients, c)
	}

	clients[0].CreateGame(1)

	for i := 1; i < NUM_CLIENT; i++ {
		clients[i].JoinGame(1)
	}

	clients[0].PlayCard()
	for {

	}

	//time.Sleep(time.Second * 10)
	//clients[0].PlayCard()
}

type Client struct {
	*websocket.Conn
	Id uint32 `json:"player_id,omitempty"`
}

func (c *Client) Start() {
	c.Conn = newWsConn()

	c.WriteJSON(&c)

	// Receive messages
	go func() {
		defer func() {
			err := c.Close()
			if err != nil {
				log.Printf("%v", err)
			}
		}()
		for {
			var state simulation.GameState
			err := c.ReadJSON(&state)
			if err != nil {
				log.Print("Error reading JSON")
			}
			log.Printf("Client %v:Message received: %v", c.Id, state)
			bytes, err := json.Marshal(&state)
			log.Printf("%s", bytes)
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
	u := url.URL{Scheme: "ws", Host: "15.188.14.213:11050", Path: "/simulation"}

	// Establish connection
	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

	return c
}
