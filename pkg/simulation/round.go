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
func (r *round) checkWinner() (bool, int) {

	if r.pos < 3 {
		return false, 0
	}
	winnerCard := r.cardsPlayed[0]
	winnerPos := 0
	for i, c := range r.cardsPlayed {

		// If they same suits and not wins
		if winnerCard.SameSuit(c) && !winnerCard.Wins(c) {
			winnerCard = c
			winnerPos = i
		} else if c.IsTriumph(r.triumph) && !winnerCard.SameSuit(c) {
			winnerCard = c
			winnerPos = i
		}
	}
	r.winner = winnerCard
	return true, winnerPos
}

func (r *round) Points() int {
	return r.points
}
