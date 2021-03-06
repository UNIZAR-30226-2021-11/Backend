package state

import (
	"fmt"
	"math/rand"
)

type Player struct {
	Cards      [6]*Card `json:"cards"`
	cardCount  int
	InternPair uint32
	Id         uint32 `json:"id"`
	Pair       uint32 `json:"pair"`
	UserName   string `json:"username"`

	CanPlay   bool   `json:"can_play"`
	CanSing   bool   `json:"can_sing"`
	SingSuit  string `json:"sing_suit"`
	CanChange bool   `json:"can_change"`

	SingingSuits []string `json:"singing_suits"`
}

// CreatePlayer Creates a new player with its ID and a pair ID.
func CreatePlayer(id uint32, pair uint32, username string) *Player {
	return &Player{
		Id:       id,
		Pair:     pair,
		UserName: username,
	}
}

func (p *Player) PlayCard(card *Card) {
	for i, c := range p.Cards {
		if c != nil && c.Equals(card) {
			p.playCard(i)
			return
		}
	}
}

// Plays the Card and compacts it.
func (p *Player) playCard(cardNumber int) {
	p.Cards[cardNumber] = nil
	for i := cardNumber; i < len(p.Cards)-1; i++ {
		p.Cards[i] = p.Cards[i+1]
	}
	p.Cards[5] = nil
	p.cardCount--
}

// Deals a Card to the Player.
func (p *Player) dealCard(card *Card) {
	p.Cards[5] = card
	p.cardCount++
}

// sameId check if is the Player with this id.
func (p *Player) sameId(id uint32) bool {
	return p.Id == id
}

// DealCards Deals 6 cards to the player.
func (p *Player) DealCards(cards [6]*Card) {
	for i, card := range cards {
		p.Cards[i] = card
	}
	p.cardCount = 6
}

// HasSing Checks if the player has singing pair and updates the record.
func (p *Player) HasSing() ([]string, bool) {
	sings := make(map[string]bool)
	var suitSings []string
	has10, has12 := false, false
	// If it doesn't have a 10 and a 12 card, exit
	for _, c := range p.Cards {
		if c != nil && c.Val == 10 {
			has10 = true
			break
		}
	}
	if !has10 {
		return nil, false
	}
	for _, c := range p.Cards {
		if c != nil && c.Val == 12 {
			has12 = true
			break
		}
	}
	if !has12 {
		return nil, false
	}

	// Check 10 && 12 of the same suit
	hasSing := false
	for i, c1 := range p.Cards {

		for j := i + 1; j < len(p.Cards); j++ {

			if c1.IsSingingPair(p.Cards[j]) {
				sings[c1.Suit] = true
				hasSing = true
			}
		}
	}
	for k := range sings {
		suitSings = append(suitSings, k)
	}

	p.SingingSuits = suitSings
	return suitSings, hasSing
}

// ChangeCard changes the seven for the top card in the deck.
func (p *Player) ChangeCard(triumph string, card *Card) {
	for _, c := range p.Cards {
		if c != nil && c.IsTriumph(triumph) && c.Val == 7 {
			c = card
		}
	}
}

// GetSeven returns the seven from the player hand.
func (p *Player) GetSeven(triumph string) *Card {
	for _, card := range p.Cards {
		if card != nil && card.IsTriumph(triumph) && card.Val == 7 {
			return card
		}
	}
	return nil
}

// PickRandomCard returns a random card from the player's hand.
func (p *Player) PickRandomCard(seed int64) (c *Card) {
	rand.Seed(seed)
	for c == nil {
		c = p.Cards[rand.Intn(p.cardCount-1)]
	}
	return c
}

// PickCard returns a card from the player's hand.
func (p *Player) PickCard(card int) (c *Card) {
	return p.Cards[card]
}

// GetCards return the non-nil cards
func (p *Player) GetCards() []*Card {
	var cs []*Card
	for _, c := range p.Cards {
		if c != nil {
			cs = append(cs, c)
		}
	}

	return cs
}

// SetPlay changes if this player can play.
func (p *Player) SetPlay(canPlay bool) {
	p.CanPlay = canPlay
}

func (p *Player) String() string {
	return fmt.Sprintf(
		"P:%v,ID:%d Pair %d",
		p.CanPlay, p.Id, p.Pair)
}
