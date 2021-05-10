package main

import (
	"Backend/pkg/events"
	"Backend/pkg/simulation"
	"Backend/pkg/state"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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
	//var controlledAI []chan events.Event
	// Creamos las ia
	enemy1, _, _ := createAI(2)
	go enemy1.controlAI()
	ally1, _, _ := createAI(1)
	go ally1.controlAI()
	enemy2, _, _ := createAI(2)
	go enemy2.controlAI()

	ps = append(ps, player, enemy1.Player, ally1.Player, enemy2.Player)

	//controlledAI = append(controlledAI, e1chan, a1chan, e2chan)

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
