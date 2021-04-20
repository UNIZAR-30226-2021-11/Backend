package simulation

import "Backend/pkg/state"

type round struct {
	pWinner     int
	sp          int
	winner      *state.Card
	suit        string
	triumph     string
	cardsPlayed [4]*state.Card
	// nº of cards played
	pos    int
	points int
}

// NewRound creates new round, receives
func NewRound(startingPair int, triumph string) *round {
	return &round{
		triumph: triumph,
		sp:      startingPair,
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
			r.updatePair(i)
		} else if c.IsTriumph(r.triumph) && !winnerCard.SameSuit(c) {
			winnerCard = c
			winnerPos = i
			r.updatePair(i)
		}
	}
	r.winner = winnerCard
	return true, winnerPos
}

// Updates the winner pair
func (r *round) updatePair(id int) {
	if id%2 == 0 {
		r.pWinner = r.sp
	} else {
		if r.sp == TeamA {
			r.pWinner = TeamB
		} else {
			r.pWinner = TeamA
		}
	}
}

func (r *round) Points() int {
	return r.points
}
