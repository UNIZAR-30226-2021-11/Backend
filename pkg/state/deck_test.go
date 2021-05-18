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
			//t.Logf("%v de %s, vale %d", c, suit, d.cards[i].Points)
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
				Val:    c,
			}
			if !carta.Equals(d.cards[i]) {
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
			if !CreateCard(suit, c).Equals(d.PickCard()) {

				t.Error("not the same card")
			}
			i++

		}
	}
	t.Run("deal 40 cards", func(t *testing.T) {

		if i != 40 {
			t.Errorf("got %v, want %v", i, 40)
		}
	})

	t.Run("top updates correctly", func(t *testing.T) {

		if d.top != 39 {
			t.Errorf("got %v, want %v", d.top, 39)
		}
	})
}

func TestChangeSeven(t *testing.T) {
	d := NewDeck()
	d.GetTriumph()
}

func TestSumCards(t *testing.T) {
	d := NewDeck()

	t.Run("sum without shuffle", func(t *testing.T) {
		sum := 0
		for _, c := range d.cards {
			sum += c.Points
		}
		if sum != 120 {
			t.Errorf("got %v, want %v", sum, 120)
		}
	})

	t.Run("sum after shuffle", func(t *testing.T) {
		d.Shuffle()
		sum := 0
		for _, c := range d.cards {
			sum += c.Points
		}
		if sum != 120 {
			t.Errorf("got %v, want %v", sum, 120)
		}
	})
}

func TestDeck_Pick4Cards(t *testing.T) {
	d := NewDeck()
	cs := d.Pick4Cards()

	t.Run("check distinct cards", func(t *testing.T) {
		for i, c := range cs {
			for j := range cs {
				if i != j {
					if c.Equals(cs[j]) {
						t.Errorf("equal cards picked")
					}
				}
			}
		}
	})

	d.Shuffle()
	t.Run("check distinct cards after shuffle", func(t *testing.T) {
		for i, c := range cs {
			for j := range cs {
				if i != j {
					if c.Equals(cs[j]) {
						t.Errorf("equal cards picked")
					}
				}
			}
		}
	})
}
