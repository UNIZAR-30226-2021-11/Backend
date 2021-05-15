package events

type UserLeft struct {
	PlayerID uint32
	GameID   uint32
	PairID   uint32
}
