package state

import (
	"encoding/json"
	"testing"
)

var (
	ID_NEW_PLAYER = 0
	PAIR = 0
)

func TestCreatePlayer(t *testing.T) {

	id := 2134
	pair := 2
	p1 := CreatePlayer(id, pair)

	if p1.Id != id {
		t.Errorf("got %v, want %v", p1.Id, id)
	}
}

func TestSerializePlayer(t *testing.T) {
	id := 2134
	pair := 2
	p1 := CreatePlayer(id, pair)
	cards := []*Card{
		CreateCard(SUIT4, 5),
		CreateCard(SUIT4, 1),
		CreateCard(SUIT1, 4),
		CreateCard(SUIT1, 11),
		CreateCard(SUIT3, 11),
		CreateCard(SUIT2, 11),
	}

	p1.DealCards(cards)
	b, err := json.Marshal(p1)
	if err != nil {
		t.Errorf("error marshaling")
	}

	t.Log(string(b))
}

func CreateTestPlayer() *Player {

	defer func() {
		ID_NEW_PLAYER++
		PAIR++
	}()

	return CreatePlayer(ID_NEW_PLAYER, PAIR % 2)
}