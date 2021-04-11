package server

import (
	"Backend/internal/data"
	"Backend/pkg/events"
	"log"
)

type Sender struct {
	sr 					*SimulationRouter
	userNameRegistry 	*data.UserNamesRegistry
}

func NewSender(sr *SimulationRouter, userNameRegistry *data.UserNamesRegistry) *Sender {
	return &Sender{
		sr:             	sr,
		userNameRegistry:   userNameRegistry,
	}
}

func (sender *Sender) HandleUserConnected(userConnectedEvent *events.UserConnected) {
	//sender.sendConstantMessage(userConnectedEvent.ClientID)
	log.Println("User connected!")
	sender.sr.SendToClient(userConnectedEvent.ClientID, "hello")
}
