package events

import "Backend/pkg/state"

type CardPlayed struct {
	PlayerID uint32
	GameID   uint32
	Card     *state.Card
}
