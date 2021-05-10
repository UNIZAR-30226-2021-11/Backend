package main

import (
	"Backend/pkg/events"
	"Backend/pkg/simulation"
	"Backend/pkg/state"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	IdNewPlayer uint32
	PAIR        uint32
	USERNAME    = "PEPE"
)

var upgrader = websocket.Upgrader{} // use default options
var addr = flag.String("a", "localhost:9000", "http service address")

func main() {
	http.HandleFunc("/simulation", mockUpGame)
	log.Fatal(http.ListenAndServe(*addr, nil))

}

func mockUpGame(w http.ResponseWriter, r *http.Request) {
	evts := make(chan events.Event)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error:", err)
		return
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Printf("closing err: %v", err)
		}
	}(c)

	// Leemos el primer evento
	var event events.GameCreate
	err = c.ReadJSON(&event)
	if err != nil {
		log.Printf("not a game create event")
		panic(err)
	}

	player := state.CreatePlayer(event.PlayerID, event.GameID, event.UserName)
	player.Pair = 1
	var ps []*state.Player
	var controlledAI []chan events.Event
	// Creamos las ia
	enemy1, e1chan := createAI(2)
	go controlAI(enemy1, e1chan, nil)
	ally1, a1chan := createAI(1)
	go controlAI(ally1, a1chan, nil)
	enemy2, e2chan := createAI(2)
	go controlAI(enemy2, e2chan, nil)

	ps = append(ps, player, enemy1, ally1, enemy2)

	controlledAI = append(controlledAI, e1chan, a1chan, e2chan)

	// TODO crear jugadores
	g := simulation.NewGame(ps)

	log.Printf("g:%v", g)
	// Read events from the websocket
	go func() {
		var evt events.Event
		for {
			err = c.ReadJSON(&evt)
			if err != nil {
				log.Println("read:", err)
				panic(err)
			}
			evts <- evt
		}
	}()

}

func controlAI(ai *state.Player, out chan events.Event, in chan *simulation.GameState) {

	for {
		select {
		case newState := <-in:
			if aiCanPlay(ai, newState) {
				// TODO react to player
			}
			if aiCanSing(ai, newState) {

			}
			if aiCanChange(ai, newState) {

			}
		}
	}
}

func aiCanPlay(ai *state.Player, newState *simulation.GameState) bool {
	for _, p := range newState.Players.All {
		if p.Id == ai.Id && p.CanPlay {
			return true
		}
	}
	return false
}

func aiCanSing(ai *state.Player, newState *simulation.GameState) bool {
	for _, p := range newState.Players.All {
		if p.Id == ai.Id && p.CanSing {
			return true
		}
	}
	return false
}

func aiCanChange(ai *state.Player, newState *simulation.GameState) bool {
	for _, p := range newState.Players.All {
		if p.Id == ai.Id && p.CanChange {
			return true
		}
	}
	return false
}

func createAI(pair uint32) (*state.Player, chan events.Event) {
	evtChan := make(chan events.Event)
	defer func() {
		IdNewPlayer++
		PAIR++
	}()
	return state.CreatePlayer(
		IdNewPlayer, pair,
		fmt.Sprintf("IA_%d", IdNewPlayer)), evtChan
}
