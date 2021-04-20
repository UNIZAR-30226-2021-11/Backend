package state

import (
	"encoding/json"
	"testing"
)

var (
	ID_NEW_PLAYER = 0
	PAIR          = 0
)

func TestCreatePlayer(t *testing.T) {

	id := 2134
	pair := 2
	p1 := CreatePlayer(uint32(id), pair)

	if p1.Id != uint32(id) {
		t.Errorf("got %v, want %v", p1.Id, id)
	}
}

func TestSerializePlayer(t *testing.T) {
	id := 2134
	pair := 2
	p1 := CreatePlayer(uint32(id), pair)
	cards := [6]*Card{
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

func TestPlayer_PlayCard(t *testing.T) {
	id := 2134
	pair := 2
	p1 := CreatePlayer(uint32(id), pair)
	cs := [6]*Card{
		CreateCard(SUIT4, 5),
		CreateCard(SUIT4, 1),
		CreateCard(SUIT1, 4),
		CreateCard(SUIT1, 11),
		CreateCard(SUIT3, 11),
		CreateCard(SUIT2, 11),
	}

	p1.DealCards(cs)

	t.Run("check full hand", func(t *testing.T) {
		if p1.cardCount != 6 {
			t.Errorf("got %v, want %v", p1.cardCount, 0)
		}
	})

	t.Run("play all cards", func(t *testing.T) {
		for _, c := range cs {
			p1.PlayCard(c)
		}
	})

	t.Run("check empty hand", func(t *testing.T) {
		if p1.cardCount > 0 {
			t.Errorf("got %v, want %v", p1.cardCount, 0)
		}
	})
}

func CreateTestPlayer() *Player {

	defer func() {
		ID_NEW_PLAYER++
		PAIR++
	}()

	return CreatePlayer(uint32(ID_NEW_PLAYER), PAIR%2)
}
