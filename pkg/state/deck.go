package state

import (
	"math/rand"
	"time"
)

const (
	SUIT1 = "oros"
	SUIT2 = "copas"
	SUIT3 = "espadas"
	SUIT4 = "bastos"
)

var (
	suits = [4]string{SUIT1, SUIT2, SUIT3, SUIT4}
	cards = [10]int{1, 2, 3, 4, 5, 6, 7, 10, 11, 12}
)

type Deck struct {
	cards   [40]*Card
	top     int
	triumph string
}

// Creates a new ordered Deck
func NewDeck() *Deck {
	baraja := Deck{}
	i := 0
	for _, suit := range suits {
		for _, c := range cards {
			baraja.cards[i] = CreateCard(suit, c)
			i++
		}
	}
	return &baraja
}

// Shuffles the Deck
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
	d.triumph = d.cards[39].Suit
}

func (d *Deck) ChangeCard(seven *Card) (c *Card) {
	d.cards[39], c = c, d.cards[39]
	return c
}

// GetTriumph returns the current triumph of the deck
func (d *Deck) GetTriumph() string {
	return d.triumph
}

// Deals the next card of the deck
func (d *Deck) PickCard() (c *Card) {
	c = d.cards[d.top]
	if d.top < 39 {

		d.top++
	}
	return c
}

func (d *Deck) Pick4Cards() (cards [4]*Card) {
	for i := range cards {
		cards[i] = d.PickCard()
	}
	return cards
}

func (d *Deck) Pick6Cards() (cards [6]*Card) {
	for i := range cards {
		cards[i] = d.PickCard()
	}
	return cards
}

func (d *Deck) InitialPick() (cards [4][6]*Card) {
	for i := range cards {
		cards[i] = d.Pick6Cards()
	}
	return cards
}
