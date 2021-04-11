package events

type UserJoined struct {
	ClientID uint32
	GameID   uint32
	UserName string
}
