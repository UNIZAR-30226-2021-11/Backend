package events

type VotePause struct {
	PlayerID uint32
	GameID   uint32
	Vote     bool
}
