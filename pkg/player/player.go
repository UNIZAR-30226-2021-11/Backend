package player

type Player struct {
	ID     		 uint      	`json:"id,omitempty"`
	UserID		 uint 		`json:"user_id,omitempty"`
	PairID 		 uint		`json:"pair_id,omitempty"`
}