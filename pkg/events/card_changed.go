package events

type CardChanged struct {
	PlayerID uint32
	GameID   uint32
	Changed  bool
}
