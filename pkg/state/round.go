package state

type round struct {
	winner 	string
	suit	string
	lastCard *card
	cardPlayed chan card
	triumph string
	//
	cardsPlayed [4]*card
}



func (r *round) playedCard(pos int, c *card)  {
	r.cardsPlayed[pos] = c
	if pos == 4 {
		r.checkWinner()
	}
}

func (r *round) checkWinner() {

	winnerCard := r.cardsPlayed[0]
	for i, c := range r.cardsPlayed{
		if i == 0 {
			continue
		}
		// If they same suits and not wins
		if winnerCard.sameSuit(c) && ! winnerCard.wins(c) {
			winnerCard = c
		}
		if c.suit == r.triumph && ! winnerCard.sameSuit(c){
			winnerCard = c
		}
	}
}
