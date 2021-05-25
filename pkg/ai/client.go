package ai

import (
	"Backend/pkg/events"
	"Backend/pkg/simulation"
	"Backend/pkg/state"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

type Client struct {
	*websocket.Conn
	P        *state.Player
	Id       uint32 `json:"player_id,omitempty"`
	PairId   uint32 `json:"pair_id,omitempty"`
	gameId   uint32
	GameData *GameData
}
type GameData struct {
	Status string               `json:"status,omitempty"`
	Game   simulation.GameState `json:"game_state,omitempty"`
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

// Start Creates a new conn a proceeds with the full protocol spec
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

// TakeOver establish a WS conn and start handling events instead of the player
func (c *Client) TakeOver() {
	log.Printf("ai %s taking over player %d", c.P.UserName, c.Id)
	c.Conn = newWsConn()
	err := c.WriteJSON(&c)
	if err != nil {
		log.Printf("ai error sending JSON:%v", err)
		return
	}
	go c.handleEvents()
}

func (c *Client) handleEvents() {
	defer func() {
		err := c.Conn.Close()
		if err != nil {
			log.Printf("ai error while closing WS:%v", err)
		}
	}()
	for {
		err := c.ReadJSON(&c.GameData)

		if err != nil {
			log.Print("ai error reading JSON")
			err := c.Conn.Close()
			if err != nil {
				return
			}
			continue
		}
		time.Sleep(time.Second * 3)
		switch c.GameData.Status {
		case "votePause":
			c.VotePause()
		case "paused":
			continue
		case "normal":
			if c.GameData.Game.Ended {
				//log.Printf("ai %d, game %d: leaving", c.Id, c.gameId)
				return
			}
			if c.CanPlay() {
				//log.Printf("ai %d, game %d: playing card", c.Id, c.gameId)
				c.PlayCard()
			}
			ok, suit := c.CanSing()
			if ok {
				//log.Printf("ai %d, game %d: singing", c.Id, c.gameId)
				c.Sing(suit)
			}

			if c.CanChange() {
				//log.Printf("ai %d, game %d: changing", c.Id, c.gameId)
				c.ChangeCard()
			}
		}

		//b, err := json.Marshal(c.GameData)
		//log.Printf("Client %v:Message received: %s", c.Id, b)
	}
}

func (c *Client) JoinGame(game uint32) {
	//log.Printf("ai %s joining game %d", c.P.UserName, game)
	event := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		PairID:    c.PairId,
		EventType: events.USER_JOINED,
	}
	err := c.WriteJSON(event)
	if err != nil {
		log.Printf("error joining game, %v", err)
	}
}

func (c *Client) CreateGame(game uint32) {
	event := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		PairID:    c.PairId,
		EventType: events.GAME_CREATE,
	}
	_ = c.WriteJSON(event)
}

func (c *Client) PlayCard() {
	cr := c.pickBestCard()
	e := events.Event{
		GameID:    c.gameId,
		PlayerID:  c.P.Id,
		EventType: events.CARD_PLAYED,
		Card:      cr,
	}
	_ = c.WriteJSON(e)
}

func (c *Client) Sing(suit string) {
	e := events.Event{
		GameID:    c.gameId,
		PlayerID:  c.Id,
		EventType: events.SING,
		Suit:      suit,
		HasSinged: true,
	}
	_ = c.WriteJSON(e)
}

func (c *Client) ChangeCard() {
	event := events.Event{
		GameID:    c.gameId,
		PlayerID:  c.Id,
		EventType: events.CARD_CHANGED,
		Changed:   true,
	}
	_ = c.WriteJSON(event)
}

func (c *Client) PauseGame(game uint32) {
	event := events.Event{
		GameID:    game,
		PlayerID:  c.Id,
		EventType: events.GAME_PAUSE,
	}
	_ = c.WriteJSON(event)
}

func (c *Client) VotePause() {
	event := events.Event{
		GameID:    c.gameId,
		PlayerID:  c.Id,
		Vote:      true,
		EventType: events.VOTE_PAUSE,
	}
	_ = c.WriteJSON(event)
}

func (c *Client) LeaveGame() {
	event := events.Event{
		GameID:    c.gameId,
		PlayerID:  c.Id,
		EventType: events.USER_LEFT,
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

func (c *Client) CanSing() (bool, string) {
	for _, p := range c.GameData.Game.Players.All {
		if p.Id == c.Id && p.CanSing {
			return true, p.SingSuit
		}
	}
	return false, ""
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
	//godotenv.Load(".env")
	//port := os.Getenv("PORT")
	//host := os.Getenv("HOST")
	u := url.URL{Scheme: "ws", Host: ":11050", Path: "/simulation"}

	// Establish connection
	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

	return c
}
