package simulation

import (
	"Backend/pkg/state"
	"testing"
)

func TestCantar(t *testing.T) {
	
}

func TestName(t *testing.T) {
	
}
func TestInitGame(t *testing.T) {

}

func createPlayer(id int, pair int, cards []*state.Card) *state.Player {
	p := state.CreatePlayer(id, pair)
	p.DealCards(cards)
	return p
}

func createGame(players []*state.Player, triumph string) *Game {
	return InitGame(players, triumph)
}
