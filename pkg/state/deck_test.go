package state

import "testing"

func TestCreateDeck(t *testing.T) {
	d := NewDeck()

	suits = [4]string{"oros", "copas", "espadas", "bastos"}
	cards = [10]int{1, 2, 3, 4, 5, 6, 7, 10, 11, 12}
	i := 0
	for _, suit := range suits {
		for _, c := range cards {

			if d.cards[i].Suit != suit {
				t.Errorf("got %v, want %v", d.cards[i].Suit, suit)
			}
			if d.cards[i].Points != getPoints(c) {
				t.Errorf("got %v, want %v", d.cards[i].Points, getPoints(c))
			}
			t.Logf("%v de %s, vale %d", c, suit, d.cards[i].Points)
			i++

		}
	}
}

func TestDeck_Shuffle(t *testing.T) {
	d := NewDeck()
	d.Shuffle()
	suits = [4]string{"oros", "copas", "espadas", "bastos"}
	cards = [10]int{1, 2, 3, 4, 5, 6, 7, 10, 11, 12}
	i := 0
	diff := 0
	for _, suit := range suits {
		for _, c := range cards {
			carta := &Card{
				Suit:   suit,
				Points: getPoints(c),
				val:    c,
			}
			if !carta.equals(d.cards[i]) {
				diff++
			}
			i++
		}
	}
	if diff < 30 {
		t.Errorf("not random enough, %d", diff)
	}
}

func TestDealAllCards(t *testing.T) {

	d := NewDeck()

	suits = [4]string{"oros", "copas", "espadas", "bastos"}
	cards = [10]int{1, 2, 3, 4, 5, 6, 7, 10, 11, 12}
	i := 0
	for _, suit := range suits {
		for _, c := range cards {
			if !CreateCard(suit, c).equals(d.DealCard()) {

				t.Error("not the same card")
			}
			i++

		}
	}
}
