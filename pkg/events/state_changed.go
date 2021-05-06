package events

type StateChanged struct {
	ClientsID []uint32
	GameData  interface{}
}
