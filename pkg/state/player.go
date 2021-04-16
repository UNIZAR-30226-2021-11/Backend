package state

type Player struct {
	// TODO elegir otro TAD
	cards   [6]*Card
	ID      int
	pair    int
	canPlay bool
}

func CreatePlayer(id int, pair int) *Player {
	return &Player{
		ID:   id,
		pair: pair,
	}
}

// Plays the Card and compacts it
func (p *Player) playCard(cardNumber int) {
	p.cards[cardNumber] = nil
	for i := cardNumber; i < len(p.cards)-1; i++ {
		p.cards[i] = p.cards[i+1]
	}
	p.cards[5] = nil
}

// Deals a Card to the Player
func (p *Player) dealCard(card *Card) {
	p.cards[5] = card
}

// sameId check if is the Player with this ID
func (p *Player) sameId(id int) bool {
	return p.ID == id
}

func (p *Player) DealCards(cards []*Card) {
	for i, card := range cards {
		p.cards[i] = card
	}
}
