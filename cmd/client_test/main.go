package main

import (
	"Backend/internal/data"
	"Backend/pkg/events"
	"Backend/pkg/state"
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
	//test2()
	test1()
}

func test2() {
	c := Client{
		Id:     5,
		PairId: 2,
	}
	game := uint32(5000)

	c.Conn = newWsConn()

	err := c.WriteJSON(&c)
	if err != nil {
		log.Printf("error sending JSON:%v", err)
		return
	}
	c.CreateSoloGame(game)
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
				log.Printf("Error reading JSON, %v", err)
				//_ = c.Conn.Close()
				continue
			}
			//time.Sleep(time.Millisecond * 200)
			switch c.GameData.Status {
			case "vote":
				c.VotePause(game)
			case "paused":
				continue
			case "normal":
				if c.GameData.Game.Ended {
					return
				}
				if c.CanPlay() {
					log.Printf("ai %d, game %d: playing card", c.Id, 6)
					err = c.PlayCard(game)
					if err != nil {
						log.Printf("%v", err)
						return
					}
				}
				ok, suit := c.CanSing()
				if ok {
					log.Printf("ai %d, game %d: singing", c.Id, 6)
					c.Sing(suit, game)
				}

				if c.CanChange() {
					log.Printf("ai %d, game %d: changing", c.Id, 6)
					c.ChangeCard(game)
				}
			}

			b, err := json.Marshal(c.GameData)
			log.Printf("Client %v:Message received: %s", c.Id, b)
		}
	}()
	time.Sleep(time.Second * 1000)
}

func test1() {
	var clients []Client
	game := uint32(5000)
	for i := 0; i < NUM_CLIENT; i++ {
		c := Client{Id: uint32(40 + i), PairId: uint32((i % 2) + 20)}
		c.Start(game)
		clients = append(clients, c)
	}
	time.Sleep(time.Second * 1)

	clients[0].CreateGame(game)

	time.Sleep(time.Second * 1)

	for i := 1; i < NUM_CLIENT; i++ {
		clients[i].JoinGame(game)
	}

	time.Sleep(time.Second * 2)

	clients[3].PauseGame(game)

	time.Sleep(time.Second * 2)

	clients[0].VotePause(game)

	// Game is paused
	time.Sleep(time.Second * 2)
	for _, c := range clients {
		c.JoinGame(game)
		time.Sleep(time.Second * 1)
	}

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

func (c *Client) Start(game uint32) {
	c.Conn = newWsConn()

	err := c.WriteJSON(&c)
	if err != nil {
		log.Printf("error sending JSON:%v", err)
		return
	}

	// Receive messages
	go func() {
		defer func() {
			err := c.Close()
			if err != nil {
				log.Printf("%v", err)
			}
		}()
		for {
			for {
				err := c.ReadJSON(&c.GameData)
				if err != nil {
					log.Print("Error reading JSON")
					err := c.Conn.Close()
					if err != nil {
						return
					}
					continue
				}
				time.Sleep(time.Millisecond * 200)
				switch c.GameData.Status {
				case "vote":
					c.VotePause(game)
				case "paused":
					continue
				case "normal":
					if c.GameData.Game.Ended {
						return
					}
					if c.CanPlay() {
						//log.Printf("ai %d, game %d: playing card", c.Id, c.gameId)
						c.PlayCard(game)
					}
					ok, suit := c.CanSing()
					if ok {
						log.Printf("ai %d, game %d: singing", c.Id, game)
						c.Sing(suit, game)
					}

					if c.CanChange() {
						log.Printf("ai %d, game %d: changing", c.Id, game)
						c.ChangeCard(game)
					}
				}

				//b, err := json.Marshal(c.GameData)
				//log.Printf("Client %v:Message received: %s", c.Id, b)
			}
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

func (c *Client) CreateSoloGame(game uint32) {
	event := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		PairID:    c.PairId,
		EventType: 8,
	}
	_ = c.WriteJSON(event)
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

func (c *Client) CanPlay() bool {
	for _, p := range c.GameData.Game.Players.All {
		if p.Id == c.Id && p.CanPlay {
			return true
		}
	}
	return false
}
func (c *Client) PlayCard(game uint32) error {
	cr := c.pickBestCard()
	e := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		EventType: events.CARD_PLAYED,
		Card:      cr,
	}
	return c.WriteJSON(e)
}
func (c *Client) CanSing() (bool, string) {
	for _, p := range c.GameData.Game.Players.All {
		if p.Id == c.Id && p.CanSing {
			return true, p.SingSuit
		}
	}
	return false, ""
}

func (c *Client) Sing(suit string, game uint32) {
	e := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		EventType: events.SING,
		Suit:      suit,
		HasSinged: true,
	}
	_ = c.WriteJSON(e)
}

func (c *Client) ChangeCard(game uint32) {
	event := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		EventType: events.CARD_CHANGED,
		Changed:   true,
	}
	_ = c.WriteJSON(event)
}
func (c *Client) CanChange() bool {
	for _, p := range c.GameData.Game.Players.All {
		if p.Id == c.Id && p.CanChange {
			return true
		}
	}
	return false
}

// pickBestCard returns the first card playable
func (c *Client) pickBestCard() *state.Card {
	// TODO picking card logic
	for _, c := range c.GetCards() {
		if c != nil && c.CanBePlayed() {
			return c
		}
	}
	return nil
}

// canKill checks if the user can win the round
func (c *Client) canKill() bool {
	return false
}

func (c *Client) GetCards() []*state.Card {
	for _, p := range c.GameData.Game.Players.All {
		if p.Id == c.Id {
			return p.Cards[:]
		}
	}
	return nil
}

func (c *Client) getKillingCards() []*state.Card {
	//TODO implement
	return nil
}

// TODO implement
func (c *Client) getCurrentWinning() *state.Card {
	gd := c.GameData.Game.CardsPlayedRound
	//winnerCard := gd[0]
	for _, c := range gd {
		if c == nil {
			continue
		}
		//if winnerCard.SameSuit() {
		//
		//}
	}
	return nil
}

func (c *Client) roundPoints() int {
	cardsplayed := c.GameData.Game.CardsPlayedRound
	points := 0
	for _, c := range cardsplayed {
		points += c.Points
	}
	return points
}

func newWsConn() *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: ":9000", Path: "/simulation"}

	// Establish connection
	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

	return c
}
