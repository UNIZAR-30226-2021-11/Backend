package game

import (
	"Backend/pkg/pair"
	"time"
)

// Game created by a user.
type Game struct {
<<<<<<< Updated upstream
	ID           uint        `json:"id,omitempty"`
	MyPairID     uint        `json:"my_pair_id,omitempty"`
	MyPlayerID   uint        `json:"my_player_id,omitempty"`
	Name         string      `json:"name,omitempty"`
	Public       bool        `json:"public,omitempty"`
	Tournament   bool        `json:"tournament,omitempty"`
	PlayersCount int         `json:"players_count,omitempty"`
	Winned       bool        `json:"winned"`
	Points       int         `json:"points"`
	Pairs        []pair.Pair `json:"pairs,omitempty"`
	CreationDate time.Time   `json:"creation_date,omitempty"`
	EndDate      time.Time   `json:"end_date,omitempty"`
}
=======
	ID     		 uint      	`json:"id,omitempty"`
	MyPairID	 uint		`json:"my_pair_id,omitempty"`
	MyPlayerID	 uint		`json:"my_player_id,omitempty"`
	Name 		 string		`json:"name,omitempty"`
	Public   	 bool		`json:"public,omitempty"`
	PlayersCount int 		`json:"players_count,omitempty"`
	Winned 		 bool		`json:"winned,omitempty"`
	WinnedPair	 uint		`json:"winned_pair,omitempty"`
	Points		 int 		`json:"points,omitempty"`
	Pairs		 []pair.Pair`json:"pairs,omitempty"`
	CreationDate time.Time 	`json:"creation_date,omitempty"`
	EndDate		 time.Time 	`json:"end_date,omitempty"`
}
>>>>>>> Stashed changes
