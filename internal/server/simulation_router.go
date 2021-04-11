package server

import (
	"Backend/internal/data"
	"Backend/pkg/events"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type SimulationRouter struct {

	clients          map[uint32]*Client
	EventsDispatcher *events.EventDispatcher
	IdManager        *data.IdManager
	userNameRegistry *data.UserNamesRegistry
	upgrader         *websocket.Upgrader

}

func NewSimulationRouter() *SimulationRouter {
	eventDispatcher := events.NewEventDispatcher()
	userNameRegistry := data.NewUserNamesRegistry()
	idManager := data.NewIdManager()

	sr := &SimulationRouter{
		clients:          make(map[uint32]*Client),
		EventsDispatcher: eventDispatcher,
		IdManager:        idManager,
		userNameRegistry: userNameRegistry,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			}},
	}

	//updater := simulation.NewUpdater()
	sender := NewSender(sr, userNameRegistry)
	eventDispatcher.RegisterUserConnectedListener(sender)

	go eventDispatcher.RunEventLoop()

	return sr
}

func (sr *SimulationRouter) Handler (w http.ResponseWriter, r *http.Request) {
	conn, err := sr.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, sr)
	sr.clients[client.ID] = client

	sr.EventsDispatcher.FireUserConnected(&events.UserConnected{ClientID: client.ID})

	log.Println("Added new client. Now", len(sr.clients), "clients connected.")
	client.Listen()
}

func (sr *SimulationRouter) SendToClient(clientID uint32, data interface{}) {
	j, err := json.Marshal(data)
	if err != nil {
		return
	}

	client, ok := sr.clients[clientID]
	if ok {
		client.SendMessage(&j)
	} else {
		log.Printf("Client %d not found\n", clientID)
		return
	}
}
