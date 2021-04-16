package events

const (
	GAME_CREATE     = 0
	USER_JOINED 	= 1
	USER_LEFT		= 2
	CARD_PLAYED 	= 3
	CARD_CHANGED    = 4
	SING 			= 5
)


// Event is a generic event communication
type Event struct {
	GameID			uint32 	`json:"game_id,omitempty"`
	PlayerID		uint32 	`json:"player_id,omitempty"`
	EventType 		int		`json:"event_type,omitempty"`
	//TODO: add extra fields when needed
}
