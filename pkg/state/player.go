package state

type player struct {
	// TODO elegir otro TAD
	cards [6]*card
}

// Plays the card and compacts it
func (p *player) playCard(id int) {
	p.cards[id] = nil
	for i := id; i < len(p.cards)-1; i++{
		p.cards[i] = p.cards[i +1]
	}
	p.cards[5] = nil
}

// Deals a card to the player
func (p *player) dealCard(card *card) {
	p.cards[5] = card
}

