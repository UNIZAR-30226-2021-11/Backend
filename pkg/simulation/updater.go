package simulation

import (
	"Backend/pkg/events"
)

type Updater struct {
	games map[uint32]*game
}

func NewUpdater() *Updater{
	return &Updater {

	}
}

func (updater *Updater) HandleNewGame() {

}

func (updater *Updater) HandleUserJoined(userJoinedEvent *events.UserJoined) {

}

func (updater *Updater) HandleUserLeft(userLeftEvent *events.UserLeft) {

}