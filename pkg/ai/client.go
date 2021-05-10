package ai

import (
	"Backend/internal/data"
	"Backend/pkg/events"
	"Backend/pkg/state"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
)

type Client struct {
	*websocket.Conn
	P        *state.Player
	Id       uint32 `json:"player_id,omitempty"`
	PairId   uint32 `json:"pair_id,omitempty"`
	gameId   uint32
	GameData *data.GameData
}

func Create(id, pairId, gameId uint32) *Client {
	p := state.CreatePlayer(id, pairId, fmt.Sprintf("IA_%d", id))
	c := Client{
		Id:     id,
		PairId: pairId,
		gameId: gameId,
		P:      p,
	}

	return &c
}

func (c *Client) Close() {
	c.Conn.Close()
}

func (c *Client) Start() {
	c.Conn = newWsConn()

	err := c.WriteJSON(&c)
	if err != nil {
		log.Printf("error sending JSON:%v", err)
		return
	}
	c.JoinGame(c.gameId)
	// Receive messages
	go c.handleEvents()
}

func (c *Client) handleEvents() {
	defer func() {
		err := c.Conn.Close()
		if err != nil {
			log.Printf("%v", err)
		}
	}()
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
		// VOTE PAUSE
		//c.PlayCard()
		bytes, err := json.Marshal(c.GameData)
		log.Printf("Client %v:Message received: %s", c.Id, bytes)
	}
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

func (c *Client) CanSing() bool {
	for _, p := range c.GameData.Game.Players.All {
		if p.Id == c.Id && p.CanSing {
			return true
		}
	}
	return false
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
	for _, card := range c.GetCards() {
		if card.CanBePlayed() {
			return card
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
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	u := url.URL{Scheme: "ws", Host: ":" + port, Path: "/simulation"}

	// Establish connection
	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

	return c
}
