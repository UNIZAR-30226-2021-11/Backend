package events

type SingleGameCreate struct {
	PlayerID uint32
	GameID   uint32
	UserName string
}
