package events

import (
	"Backend/pkg/simulation"
)

type StateChanged struct {
	ClientsID   []uint32
	Game	 	*simulation.Game
}