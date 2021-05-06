package main

import (
	"Backend/pkg/events"
	"Backend/pkg/simulation"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

const (
	NUM_CLIENT = 4
)

func main() {
	var clients []Client

	for i := 0; i < NUM_CLIENT; i++ {
		c := Client{Id: uint32(40 + i), PairId: uint32((i % 2) + 20)}
		c.Start()
		clients = append(clients, c)
	}

	clients[0].CreateGame(1)

	for i := 1; i < NUM_CLIENT; i++ {
		clients[i].JoinGame(1)
	}

	time.Sleep(time.Second * 100)

	//clients[3].PauseGame(1)
}

type Client struct {
	*websocket.Conn
	Id        uint32 `json:"player_id,omitempty"`
	PairId    uint32 `json:"pair_id,omitempty"`
	GameState *simulation.GameState
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
			err := c.ReadJSON(&c.GameState)
			if err != nil {
				log.Print("Error reading JSON")
			}
			c.PlayCard()
			bytes, err := json.Marshal(c.GameState)
			log.Printf("Client %v:Message received: %s", c.Id, bytes)
		}
	}()
}

func (c *Client) JoinGame(game uint32) {
	event := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		PairID:    c.PairId,
		EventType: 1,
	}
	_ = c.WriteJSON(event)
}

func (c *Client) CreateGame(game uint32) {
	event := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		PairID:    c.PairId,
		EventType: 0,
	}
	_ = c.WriteJSON(event)
}

func (c *Client) PlayCard() {
	players := c.GameState.Players.All
	for _, player := range players {
		if player.CanPlay {
			event := events.Event{
				GameID:    1,
				PlayerID:  c.Id,
				EventType: 3,
				Card:      player.Cards[0],
			}
			_ = c.WriteJSON(event)
		}
	}
}

func (c *Client) PauseGame(game uint32) {
	event := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
	}
	_ = c.WriteJSON(event)
}

func (c *Client) LeaveGame() {
	event := events.Event{
		GameID:    1,
		PlayerID:  c.Id,
		EventType: 2,
	}
	_ = c.WriteJSON(event)
}

func newWsConn() *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: ":9000", Path: "/simulation"}

	// Establish connection
	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

	return c
}
