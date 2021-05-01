package pair

import "Backend/pkg/user"

type Pair struct {
	ID     uint        `json:"id,omitempty"`
	GameID uint        `json:"game_id,omitempty"`
	Users  []user.User `json:"users"`
}
