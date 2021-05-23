package pair

import "Backend/pkg/user"

type Pair struct {
	ID         uint        `json:"id,omitempty"`
	GameID     uint        `json:"game_id,omitempty"`
	Winned     bool        `json:"winned,omitempty"`
	GamePoints int         `json:"game_points,omitempty"`
	Users      []user.User `json:"users,omitempty"`
}
