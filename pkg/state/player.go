package state

import (
	"math/rand"
)

type Player struct {
	Cards     [6]*Card `json:"cards"`
	cardCount int
	Id        uint32 `json:"id"`
	Pair      int    `json:"pair"`
	CanPlay   bool   `json:"can_play"`
	CanSing   bool   `json:"can_sing"`
	CanChange bool   `json:"can_change"`

	SingingSuits []string `json:"singing_suits"`
	singSuit1    bool     `json:"sing_oros"`
	singSuit2    bool     `json:"sing_copas"`
	singSuit3    bool     `json:"sing_espadas"`
	singSuit4    bool     `json:"sing_bastos"`
}

// CreatePlayer Creates a new player with its ID and a pair ID
func CreatePlayer(id uint32, pair int) *Player {
	return &Player{
		Id:   id,
		Pair: pair,
	}
}

func (p *Player) PlayCard(card *Card) {
	for i, c := range p.Cards {
		if c != nil && c.Equals(card) {
			p.playCard(i)
		}
	}
}

// Plays the Card and compacts it
func (p *Player) playCard(cardNumber int) {
	p.Cards[cardNumber] = nil
	for i := cardNumber; i < len(p.Cards)-1; i++ {
		p.Cards[i] = p.Cards[i+1]
	}
	p.Cards[5] = nil
	p.cardCount--
}

// Deals a Card to the Player
func (p *Player) dealCard(card *Card) {
	p.Cards[5] = card
	p.cardCount++
}

// sameId check if is the Player with this id
func (p *Player) sameId(id uint32) bool {
	return p.Id == id
}

// DealCards Deals 6 cards to the player
func (p *Player) DealCards(cards [6]*Card) {
	for i, card := range cards {
		p.Cards[i] = card
	}
	p.cardCount = 6
}

// HasSing Checks if the player has singing pair and updates the record
func (p *Player) HasSing() ([]string, bool) {
	var suitSings []string
	for _, c1 := range p.Cards {
		for _, c2 := range p.Cards {
			if c1.IsSingingPair(c2) {
				suitSings = append(suitSings, c1.Suit)
			}
		}
	}
	p.SingingSuits = suitSings
	return suitSings, len(suitSings) > 0
}

// ChangeCard changes the seven for the top card in the deck
func (p *Player) ChangeCard(card *Card) {

}

// GetSeven returns the seven from the player hand
func (p *Player) GetSeven(triumph string) *Card {
	for _, card := range p.Cards {
		if card.IsTriumph(triumph) && card.Val == 7 {
			return card
		}
	}
	return nil
}

// PickRandomCard returns a random card from the player's hand
func (p *Player) PickRandomCard(seed int64) (c *Card) {
	rand.Seed(seed)

	return p.Cards[rand.Intn(p.cardCount-1)]
}

// PickCard returns a card from the player's hand
func (p *Player) PickCard(card int) (c *Card) {
	return p.Cards[card]
}

func (p *Player) SetPlay(canPlay bool) {
	p.CanPlay = canPlay
}
