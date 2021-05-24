package server

import (
	"Backend/pkg/events"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	channelBufSize = 100

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 6 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Client struct holds client-specific variables.
type Client struct {
	ID     uint32 `json:"player_id,omitempty"`
	pairID uint32
	gameID uint32
	ws     *websocket.Conn
	ch     chan interface{}
	doneCh chan bool
	sr     *SimulationRouter
}

// NewClient initializes a new Client struct with given websocket.
func NewClient(ws *websocket.Conn, sr *SimulationRouter) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	ch := make(chan interface{}, channelBufSize)
	doneCh := make(chan bool, channelBufSize)

	player := struct {
		Id     uint32 `json:"player_id,omitempty"`
		PairId uint32 `json:"pair_id,omitempty"`
	}{}
	err := ws.ReadJSON(&player)
	if err != nil {
		log.Print("Error reading player ID")
	}

	return &Client{
		ID:     player.Id,
		pairID: player.PairId,
		ws:     ws,
		ch:     ch,
		doneCh: doneCh,
		sr:     sr,
	}
}

// Conn returns client's websocket.Conn struct.
func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

// SendMessage sends game state to the client.
func (c *Client) SendMessage(data interface{}) {
	select {
	case c.ch <- data:
	default:
		//c.sr.monitor.AddDroppedMessage()
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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.sr.EventsDispatcher.FireUserLeft(&events.UserLeft{
			PlayerID: c.ID,
			PairID:   c.pairID,
			GameID:   c.gameID,
		})
		ticker.Stop()
		err := c.ws.Close()
		if err != nil {
			log.Println("Error:", err.Error())
		}
	}()

	log.Println("Listening write to client")
	for {
		select {

		case data := <-c.ch:
			//before := time.Now()
			err := c.ws.WriteJSON(data)
			//after := time.Now()

			if err != nil {
				log.Println(err)
			} else {
				//elapsed := after.Sub(before)
				//c.sr.monitor.AddSendTime(elapsed)
			}

		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.doneCh:
			c.doneCh <- true
			return
		}
	}
}

func (c *Client) listenRead() {
	defer func() {
		c.sr.EventsDispatcher.FireUserLeft(&events.UserLeft{
			PlayerID: c.ID,
			PairID:   c.pairID,
			GameID:   c.gameID,
		})
		err := c.ws.Close()
		if err != nil {
			log.Println("Error:", err.Error())
		}
	}()

	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

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
	var event events.Event
	err := c.ws.ReadJSON(&event)
	if err != nil {
		c.sr.EventsDispatcher.FireUserLeft(&events.UserLeft{
			PlayerID: c.ID,
			PairID:   c.pairID,
			GameID:   c.gameID,
		})
		log.Println(err)

		c.doneCh <- true
	} else {
		c.unmarshalUserInput(event)
	}
}

func (c *Client) unmarshalUserInput(event events.Event) {
	switch event.EventType {

	case events.GAME_CREATE:
		e := &events.GameCreate{
			PlayerID: event.PlayerID,
			PairID:   event.PairID,
			GameID:   event.GameID,
			UserName: event.UserName,
		}
		c.sr.EventsDispatcher.FireGameCreate(e)

	case events.SINGLE_GAME_CREATE:
		e := &events.SingleGameCreate{
			PlayerID: event.PlayerID,
			GameID:   event.GameID,
			UserName: event.UserName,
		}
		c.sr.EventsDispatcher.FireSingleGameCreate(e)

	case events.GAME_PAUSE:
		e := &events.GamePause{
			PlayerID: event.PlayerID,
			GameID:   event.GameID,
		}
		c.sr.EventsDispatcher.FireGamePause(e)

	case events.VOTE_PAUSE:
		e := &events.VotePause{
			PlayerID: event.PlayerID,
			GameID:   event.GameID,
			Vote:     event.Vote,
		}
		c.sr.EventsDispatcher.FireVotePause(e)

	case events.USER_JOINED:
		e := &events.UserJoined{
			PlayerID: event.PlayerID,
			PairID:   event.PairID,
			GameID:   event.GameID,
			UserName: event.UserName,
		}
		c.sr.EventsDispatcher.FireUserJoined(e)

	case events.USER_LEFT:
		e := &events.UserLeft{
			PlayerID: event.PlayerID,
			GameID:   event.GameID,
			PairID:   event.PairID,
		}
		c.sr.EventsDispatcher.FireUserLeft(e)

	case events.CARD_PLAYED:
		e := &events.CardPlayed{
			PlayerID: event.PlayerID,
			GameID:   event.GameID,
			Card:     event.Card,
		}
		c.sr.EventsDispatcher.FireCardPlayed(e)

	case events.CARD_CHANGED:
		e := &events.CardChanged{
			PlayerID: event.PlayerID,
			GameID:   event.GameID,
			Changed:  event.Changed,
		}
		c.sr.EventsDispatcher.FireCardChanged(e)

	case events.SING:
		e := &events.Sing{
			PlayerID:  event.PlayerID,
			GameID:    event.GameID,
			Suit:      event.Suit,
			HasSinged: event.HasSinged,
		}
		c.sr.EventsDispatcher.FireSing(e)

	default:
		log.Fatalln("Unknown message type %T", event.EventType)
	}
}
