package simulation

// New event
type Event interface {
	GetPlayerId() int
}

type PlayerEvent struct {
	playerId int
}

func (p PlayerEvent) GetPlayerId() int {
	return p.playerId
}

type CardPlayedEvent struct {
	PlayerEvent
	suit string
	val int
}

type SingEvent struct {
	PlayerEvent
	singSuit string
}

type CardChangeEvent struct {
	PlayerEvent
}


func NewGuinote() {

}
