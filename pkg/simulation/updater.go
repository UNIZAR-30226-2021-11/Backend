package simulation

import (
	"Backend/pkg/events"
	"fmt"
)

type Updater struct {
	games map[uint32]*Game
}

func NewUpdater() *Updater{
	return &Updater {

	}
}

func (updater *Updater) HandleNewGame() {

}

func (updater *Updater) HandleUserJoined(userJoinedEvent *events.UserJoined) {
	fmt.Print("User joined!")
}

func (updater *Updater) HandleUserLeft(userLeftEvent *events.UserLeft) {

}