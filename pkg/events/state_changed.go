package events

type StateChanged struct {
	ClientsID []uint32
	GameID    uint32
	GameData  interface{}
}
