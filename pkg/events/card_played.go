package events

import "Backend/pkg/state"

type CardPlayed struct {
	ClientID uint32
	GameID   uint32
	Card	 *state.Card
}
