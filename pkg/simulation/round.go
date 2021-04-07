package simulation

import "Backend/pkg/state"

type round struct {
	winner      *state.Card
	suit        string
	triumph     string
	cardsPlayed [4]*state.Card
	// nº of cards played
	pos    int
	points int
}

func NewRound(triumph string) *round {
	return &round{
		triumph: triumph,
	}
}

func (r *round) playedCard(c *state.Card) {
	if r.pos == 0 {
		r.suit = c.Suit
	}
	r.cardsPlayed[r.pos] = c
	// Sum of points
	r.points += c.Points
	if r.pos < 3 {

		r.pos++
	}
}

// Returns true if there is a winner
func (r *round) checkWinner() bool {

	if r.pos < 3 {
		return false
	}
	winnerCard := r.cardsPlayed[0]
	for _, c := range r.cardsPlayed {

		// If they same suits and not wins
		if winnerCard.SameSuit(c) && !winnerCard.Wins(c) {
			winnerCard = c
		} else if c.IsTriumph(r.triumph) && !winnerCard.SameSuit(c) {
			winnerCard = c
		}
	}
	r.winner = winnerCard
	return true
}

func (r *round) Points() int {
	return r.points
}
