package v1

import (
	"Backend/internal/data"
	"Backend/internal/server"
	"Backend/pkg/events"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type SimulationRouter struct {

	clients          map[uint32]*server.Client
	EventsDispatcher *events.EventDispatcher
	IdManager        *data.IdManager
	upgrader         *websocket.Upgrader

}

func NewSimulationRouter(eventDispatcher *events.EventDispatcher, idManager *data.IdManager) *SimulationRouter{
	return &SimulationRouter{
		clients:          make(map[uint32]*server.Client),
		EventsDispatcher: eventDispatcher,
		IdManager:        idManager,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			}},
	}
}

func (sr *SimulationRouter) Handler (w http.ResponseWriter, r *http.Request) {
	conn, err := sr.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := server.NewClient(conn, sr)
	sr.clients[client.ID] = client

	sr.EventsDispatcher.FireUserConnected(&events.UserConnected{ClientID: client.ID})

	log.Println("Added new client. Now", len(sr.clients), "clients connected.")
	client.Listen()
}