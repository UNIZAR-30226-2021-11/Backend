package main

import (
	"Backend/internal/data"
	"Backend/pkg/events"
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
	time.Sleep(time.Second * 1)

	clients[0].CreateGame(1)

	for i := 1; i < NUM_CLIENT; i++ {
		clients[i].JoinGame(1)
	}

	time.Sleep(time.Second * 2)

	clients[3].PauseGame(1)

	time.Sleep(time.Second * 2)

	clients[0].VotePause(1)

	for {
		//Guarrisimo
		time.Sleep(time.Second * 5)
	}
}

type Client struct {
	*websocket.Conn
	Id       uint32 `json:"player_id,omitempty"`
	PairId   uint32 `json:"pair_id,omitempty"`
	GameData *data.GameData
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
			err := c.ReadJSON(&c.GameData)
			if err != nil {
				log.Print("Error reading JSON")
			}
			//c.PlayCard()
			bytes, err := json.Marshal(c.GameData)
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
	players := c.GameData.Game.Players.All

guarrada:
	for _, player := range players {
		if player.Id == c.Id {
			if player.CanPlay {
				for _, card := range player.Cards {
					if card != nil && card.Playable {
						event := events.Event{
							GameID:    1,
							PlayerID:  c.Id,
							EventType: 3,
							Card:      card,
						}
						_ = c.WriteJSON(event)
						break guarrada
					}
				}
			} else if player.CanSing {
				event := events.Event{
					GameID:    1,
					PlayerID:  c.Id,
					EventType: 5,
					Suit:      player.SingSuit,
					HasSinged: false,
				}
				_ = c.WriteJSON(event)
				break guarrada
			} else if player.CanChange {
				event := events.Event{
					GameID:    1,
					PlayerID:  c.Id,
					EventType: 4,
					Changed:   false,
				}
				_ = c.WriteJSON(event)
				break guarrada
			}
		}

	}
}

func (c *Client) PauseGame(game uint32) {
	event := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		EventType: 6,
	}
	_ = c.WriteJSON(event)
}

func (c *Client) VotePause(game uint32) {
	event := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		Vote:      true,
		EventType: 7,
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
