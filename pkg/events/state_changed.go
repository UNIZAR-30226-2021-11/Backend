package events

type StateChanged struct {
	ClientsID   []uint32
	//Game	 	*simulation.Game
	Game 		interface{}
}