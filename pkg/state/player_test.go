package state

import (
	"encoding/json"
	"testing"
)

var (
	ID_NEW_PLAYER uint32
	PAIR          uint32
	USERNAME      = "PEPE"
)

func TestCreatePlayer(t *testing.T) {

	p1 := CreatePlayer(ID_NEW_PLAYER, PAIR, USERNAME)

	if p1.Id != ID_NEW_PLAYER {
		t.Errorf("got %v, want %v", p1.Id, ID_NEW_PLAYER)
	}
}

func TestSerializePlayer(t *testing.T) {

	p1 := CreateTestPlayer()

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
	p1 := CreateTestPlayer()
	cs := [6]*Card{
		CreateCard(SUIT4, 5),
		CreateCard(SUIT4, 1),
		CreateCard(SUIT1, 4),
		CreateCard(SUIT1, 7),
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

	t.Run("change 7", func(t *testing.T) {
		p1.ChangeCard(SUIT1, CreateCard(SUIT1, 1))

		if p1.GetSeven(SUIT1) != nil {
			t.Errorf("player still has seven")
		}
	})
}

func CreateTestPlayer() *Player {

	defer func() {
		ID_NEW_PLAYER++
		PAIR++
	}()

	return CreatePlayer(ID_NEW_PLAYER, PAIR%2, USERNAME)
}
