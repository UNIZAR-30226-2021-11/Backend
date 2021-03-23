package games

import "time"

// Game created by a user.
type Game struct {
	ID     		 uint      	`json:"id,omitempty"`
	Name 		 string		`json:"name,omitempty"`
	Public   	 bool		`json:"public,omitempty"`
	CreationDate time.Time 	`json:"creation_date,omitempty"`
	EndDate		 time.Time 	`json:"end_date,omitempty"`
}