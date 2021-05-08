package events

import "Backend/pkg/state"

const (
	GAME_CREATE  = 0
	USER_JOINED  = 1
	USER_LEFT    = 2
	CARD_PLAYED  = 3
	CARD_CHANGED = 4
	SING         = 5
	GAME_PAUSE   = 6
	VOTE_PAUSE   = 7
)

// Event is a generic event communication
type Event struct {
	GameID    uint32      `json:"game_id,omitempty"`
	PlayerID  uint32      `json:"player_id,omitempty"`
	PairID    uint32      `json:"pair_id,omitempty"`
	UserName  string      `json:"username,omitempty"`
	EventType int         `json:"event_type,omitempty"`
	Card      *state.Card `json:"card,omitempty"`
	Changed   bool        `json:"changed,omitempty"`
	Suit      string      `json:"suit,omitempty"`
	HasSinged bool        `json:"has_singed,omitempty"`
	Vote      bool        `json:"vote,omitempty"`
}
