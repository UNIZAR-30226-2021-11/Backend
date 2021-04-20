package events

type CardChanged struct {
	ClientID uint32
	GameID   uint32
	Changed  bool
}
