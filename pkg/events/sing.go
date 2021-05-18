package events

type Sing struct {
	PlayerID  uint32
	GameID    uint32
	Suit      string
	HasSinged bool
}
