package game

import "time"

// Game created by a user.
type Game struct {
	ID     		 uint      	`json:"id,omitempty"`
	Name 		 string		`json:"name,omitempty"`
	Public   	 bool		`json:"public,omitempty"`
	PlayersCount int 		`json:"players_count,omitempty"`
	Winned 		 bool		`json:"winned,omitempty"`
	CreationDate time.Time 	`json:"creation_date,omitempty"`
	EndDate		 time.Time 	`json:"end_date,omitempty"`
}