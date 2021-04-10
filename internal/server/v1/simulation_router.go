package v1

import (
	"Backend/internal/server"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	""
)

type SimulationRouter struct {

	clients          map[uint32]*server.Client
	monitor          *monitor.Monitor
	eventsDispatcher *events.EventDispatcher
	idManager        *utils.IdManager
	upgrader		 *websocket.Upgrader

}

func (sr *SimulationRouter) simulationHandler (w http.ResponseWriter, r *http.Request) {
	conn, err := sr.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := server.NewClient(conn, sr)
	sr.clients[client.id] = client

	sr.eventsDispatcher.FireUserConnected(&events.UserConnected{ClientID: client.id})

	log.Println("Added new client. Now", len(sr.clients), "clients connected.")
	client.Listen()
}