package main

import (
	"Backend/pkg/events"
	"Backend/pkg/simulation"
	"Backend/pkg/state"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const STATES = 5

var (
	ID_NEW_PLAYER uint32
	PAIR          = 0
)

var testGames = createSimulationGames()

var upgrader = websocket.Upgrader{} // use default options
var addr = flag.String("a", "localhost:9000", "http service address")

type TestGames struct {
	count int
	games []*TestState
}

type TestState struct {
	Status string                `json:"status"`
	Game   *simulation.GameState `json:"game"`
}

func main() {

	http.HandleFunc("/simulation", simStates)
	log.Fatal(http.ListenAndServe(*addr, nil))

}

func simStates(w http.ResponseWriter, r *http.Request) {
	evts := make(chan events.Event)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	defer close(evts)
	// Asincronamente la respuesta
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case evt := <-evts:
				ticker.Reset(time.Second * 5)
				err = c.WriteJSON(testGames.games[testGames.count])
				if err != nil {
					log.Println("write:", err)
					break
				}
				log.Printf("Event %v", evt)
				testGames.count = (testGames.count + 1) % STATES

			case <-ticker.C:
				err = c.WriteJSON(testGames.games[testGames.count])
				if err != nil {
					log.Println("write:", err)
					break
				}
				log.Printf("Tick, sent state %v", testGames.count)
				testGames.count = (testGames.count + 1) % STATES
			}
		}
	}()

	// Recepción de mensajes
	var evt events.Event
	for {
		err := c.ReadJSON(&evt)
		if err != nil {
			log.Println("read:", err)
			break
		}
		evts <- evt
	}
}

func createSimulationGames() *TestGames {
	ps := createTestPlayers()
	tg := TestGames{}

	singing := createGameState(ps)
	singing.Players.Current().CanSing = true
	tg.games = append(tg.games, &TestState{
		Status: "singing",
		Game:   singing,
	})

	playing := createGameState(ps)
	playing.Players.Current().CanPlay = true
	for i := range playing.Players.Current().Cards {
		playing.Players.Current().Cards[i].Playable = true
	}
	tg.games = append(tg.games, &TestState{
		Status: "playing",
		Game:   playing,
	})

	changing := createGameState(ps)
	changing.Players.Current().CanChange = true
	tg.games = append(tg.games, &TestState{
		Status: "changing",
		Game:   changing,
	})

	vueltas := createGameState(ps)
	vueltas.Vueltas = true
	tg.games = append(tg.games, &TestState{
		Status: "vueltas",
		Game:   vueltas,
	})

	arrastre := createGameState(ps)
	arrastre.Arrastre = true
	tg.games = append(tg.games, &TestState{
		Status: "arrastre",
		Game:   arrastre,
	})

	return &tg
}

func createTestPlayers() []*state.Player {
	var players []*state.Player
	for i := 0; i < 4; i++ {
		players = append(players, CreateTestPlayer())
	}
	return players
}

func createGameState(players []*state.Player) *simulation.GameState {
	return &simulation.NewGame(players).GameState
}

func CreateTestPlayer() *state.Player {

	defer func() {
		ID_NEW_PLAYER++
		PAIR++
	}()

	return state.CreatePlayer(ID_NEW_PLAYER, PAIR%2+1)
}
