package events

const (
	USER_JOINED 	= 0
	USER_LEFT		= 1
)


// Event is a generic event communication
type Event struct {
	GameID			uint32 	`json:"game_id,omitempty"`
	PlayerID		uint32 	`json:"player_id,omitempty"`
	EventType 		int		`json:"event_type,omitempty"`
	//TODO: add extra fields when needed
}
