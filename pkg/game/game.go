package game

import (
	"Backend/pkg/pair"
	"time"
)

// Game created by a user.
type Game struct {
	ID           uint        `json:"id,omitempty"`
	MyPairID     uint        `json:"my_pair_id,omitempty"`
	MyPlayerID   uint        `json:"my_player_id,omitempty"`
	Name         string      `json:"name,omitempty"`
	Public       bool        `json:"public,omitempty"`
	Tournament   bool        `json:"tournament,omitempty"`
	PlayersCount int         `json:"players_count"`
	Winned       bool        `json:"winned"`
	WinnedPair   uint        `json:"winned_pair,omitempty"`
	WinnedPoints int         `json:"winned_points,omitempty"`
	LostPair     uint        `json:"lost_pair,omitempty"`
	LostPoints   int         `json:"lost_points,omitempty"`
	Points       int         `json:"points"`
	Pairs        []pair.Pair `json:"pairs,omitempty"`
	CreationDate time.Time   `json:"creation_date,omitempty"`
	EndDate      time.Time   `json:"end_date,omitempty"`
}
