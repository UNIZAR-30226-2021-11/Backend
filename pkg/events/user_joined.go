package events

type UserJoined struct {
	PlayerID uint32
	PairID	 uint32
	GameID   uint32
	UserName string
}
