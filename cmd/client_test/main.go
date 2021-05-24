package main

import (
	"Backend/internal/data"
	"Backend/pkg/events"
	pair2 "Backend/pkg/pair"
	"Backend/pkg/state"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	NUM_CLIENT = 4
)

func main() {
	//test3()
	//test2()
	test4()
}

func test3() {
	pair := pair2.Pair{
		Winned:     true,
		GamePoints: 80,
	}

	// initialize http client
	client := &http.Client{}

	// marshal User to json
	json, err := json.Marshal(pair)
	if err != nil {
		panic(err)
	}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, "http://localhost:9000/api/v1/pairs/18", bytes.NewBuffer(json))
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)
}

func test2() {
	c := Client{
		Id:     5,
		PairId: 2,
	}

	c.Conn = newWsConn()

	err := c.WriteJSON(&c)
	if err != nil {
		log.Printf("error sending JSON:%v", err)
		return
	}
	c.CreateSoloGame(6)
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
					c.VotePause(6)
				case "paused":
					continue
				case "normal":
					if c.GameData.Game.Ended {
						return
					}
					if c.CanPlay() {
						log.Printf("ai %d, game %d: playing card", c.Id, 6)
						err = c.PlayCard()
						if err != nil {
							log.Printf("%v", err)
							return
						}
					}
					ok, suit := c.CanSing()
					if ok {
						log.Printf("ai %d, game %d: singing", c.Id, 6)
						c.Sing(suit)
					}

					if c.CanChange() {
						log.Printf("ai %d, game %d: changing", c.Id, 6)
						c.ChangeCard()
					}
				}

				//b, err := json.Marshal(c.GameData)
				//log.Printf("Client %v:Message received: %s", c.Id, b)
			}
		}
	}()
	time.Sleep(time.Second * 1000)
}

func test4() {
	var clients []Client

	for i := 0; i < NUM_CLIENT; i++ {
		c := Client{Id: uint32(40 + i), PairId: uint32((i % 2) + 20)}
		c.Start()
		clients = append(clients, c)
	}
	time.Sleep(time.Second * 1)

	clients[0].CreateGame(6)

	time.Sleep(time.Second * 1)

	for i := 1; i < NUM_CLIENT; i++ {
		clients[i].JoinGame(6)
	}

	time.Sleep(time.Second * 5)
	clients[2].LeaveGame(6)
	time.Sleep(time.Second * 5)
	clients[2].JoinGame(6)
	for {
		//Guarrisimo
		time.Sleep(time.Second * 5)
		clients[2].LeaveGame(6)
		time.Sleep(time.Second * 5)
		clients[2].JoinGame(6)
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
					c.VotePause(6)
				case "paused":
					continue
				case "normal":
					if c.GameData.Game.Ended {
						log.Printf("client %d leaving game", c.Id)
						c.LeaveGame(6)
						return
					}
					if c.CanPlay() {
						//log.Printf("ai %d, game %d: playing card", c.Id, c.gameId)
						c.PlayCard()
					}
					ok, suit := c.CanSing()
					if ok {
						log.Printf("ai %d, game %d: singing", c.Id, 6)
						c.Sing(suit)
					}

					if c.CanChange() {
						log.Printf("ai %d, game %d: changing", c.Id, 6)
						c.ChangeCard()
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

func (c *Client) LeaveGame(game uint32) {
	c.Conn.Close()
	return
	event := events.Event{
		GameID:    game,
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
func (c *Client) PlayCard() error {
	cr := c.pickBestCard()
	e := events.Event{
		GameID:    6,
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

func (c *Client) Sing(suit string) {
	e := events.Event{
		GameID:    6,
		PlayerID:  c.Id,
		EventType: events.SING,
		Suit:      suit,
		HasSinged: true,
	}
	_ = c.WriteJSON(e)
}

func (c *Client) ChangeCard() {
	event := events.Event{
		GameID:    6,
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
