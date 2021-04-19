package state

type Player struct {
	Cards   [6]*Card `json:"cards"`
	Id      int `json:"id"`
	Pair    int `json:"pair"`
	CanPlay bool `json:"can_play"`
	CanSing bool `json:"can_sing"`
}

// CreatePlayer Creates a new player with its ID and a pair ID
func CreatePlayer(id int, pair int) *Player {
	return &Player{
		Id:   id,
		Pair: pair,
	}
}

// Plays the Card and compacts it
func (p *Player) playCard(cardNumber int) {
	p.Cards[cardNumber] = nil
	for i := cardNumber; i < len(p.Cards)-1; i++ {
		p.Cards[i] = p.Cards[i+1]
	}
	p.Cards[5] = nil
}

// Deals a Card to the Player
func (p *Player) dealCard(card *Card) {
	p.Cards[5] = card
}

// sameId check if is the Player with this id
func (p *Player) sameId(id int) bool {
	return p.Id == id
}

// DealCards Deals 6 cards to the player
func (p *Player) DealCards(cards []*Card) {
	for i, card := range cards {
		p.Cards[i] = card
	}
}

// HasSing Checks if the player has singing pair
func (p *Player) HasSing() ([]string, bool) {
	var suitSings []string
	for _, c1 := range p.Cards {
		for _,c2 := range p.Cards {
			if c1.IsSingingPair(c2){
				suitSings = append(suitSings, c1.Suit)
			}
		}
	}
	return suitSings, len(suitSings) > 0
}