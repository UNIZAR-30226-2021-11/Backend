package server

import (
	"Backend/internal/data"
	"Backend/pkg/events"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type SimulationRouter struct {
	clients              map[uint32]*Client
	EventsDispatcher     *events.EventDispatcher
	userNameRegistry     *data.UserNamesRegistry
	upgrader             *websocket.Upgrader
	simulationRepository *data.SimulationRepository
}

func NewSimulationRouter() *SimulationRouter {
	eventDispatcher := events.NewEventDispatcher()
	userNameRegistry := data.NewUserNamesRegistry()

	sr := &SimulationRouter{
		clients:          make(map[uint32]*Client),
		EventsDispatcher: eventDispatcher,
		userNameRegistry: userNameRegistry,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			}},
		simulationRepository: data.NewSimulationRepository(eventDispatcher),
	}

	eventDispatcher.RegisterStateChangedListener(sr)
	eventDispatcher.RegisterGameCreateListener(sr.simulationRepository)
	eventDispatcher.RegisterSingleGameCreateListener(sr.simulationRepository)
	eventDispatcher.RegisterGamePauseListener(sr.simulationRepository)
	eventDispatcher.RegisterVotePauseListener(sr.simulationRepository)
	eventDispatcher.RegisterUserJoinedListener(sr.simulationRepository)
	eventDispatcher.RegisterUserJoinedListener(sr)
	eventDispatcher.RegisterUserLeftListener(sr.simulationRepository)
	eventDispatcher.RegisterUserLeftListener(sr)
	eventDispatcher.RegisterCardPlayedListener(sr.simulationRepository)
	eventDispatcher.RegisterCardChangedListener(sr.simulationRepository)
	eventDispatcher.RegisterSingListener(sr.simulationRepository)

	go eventDispatcher.RunEventLoop()

	return sr
}

func (sr *SimulationRouter) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := sr.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, sr)
	sr.clients[client.ID] = client

	log.Println("Added new client. Now", len(sr.clients), "clients connected.")
	client.Listen()
}

func (sr *SimulationRouter) SendToClient(clientID uint32, data interface{}) {
	client, ok := sr.clients[clientID]
	if ok {
		client.SendMessage(data)
	} else {
		log.Printf("Client %d not found\n", clientID)
		return
	}
}

func (sr *SimulationRouter) HandleStateChanged(stateChangedEvent *events.StateChanged) {
	clientsID := stateChangedEvent.ClientsID
	for _, c := range clientsID {
		sr.SendToClient(c, &stateChangedEvent.GameData)
	}
}

func (sr *SimulationRouter) HandleUserLeft(userLeftEvent *events.UserLeft) {
	clientID := userLeftEvent.PlayerID
	client, ok := sr.clients[clientID]
	if !ok {
		log.Printf("Client %d not found\n", clientID)
		return
	}
	client.Done()
	delete(sr.clients, clientID)
}

func (sr *SimulationRouter) HandleUserJoined(userJoinedEvent *events.UserJoined) {
	clientID := userJoinedEvent.PlayerID
	client, ok := sr.clients[clientID]
	if !ok {
		log.Printf("Client %d not found\n", clientID)
		return
	}
	client.gameID = userJoinedEvent.GameID
}
