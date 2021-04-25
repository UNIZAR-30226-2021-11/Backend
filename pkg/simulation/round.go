package simulation

import (
	"Backend/pkg/state"
	"fmt"
)

type rounds struct {
	rondas  [10]*round
	current int
	triumph string
}

func (r *rounds) next() error {
	if r.current > 9 {
		return fmt.Errorf("max rounds passed")
	}
	r.rondas[r.current] = NewRound(r.triumph)
	return nil
}

type round struct {
	winner      *cardPlayed
	suit        string
	triumph     string
	cardsPlayed [4]*state.Card

	played [4]*cardPlayed

	// nº of cards played
	pos    int
	points int
}

type cardPlayed struct {
	*state.Player
	*state.Card
}

// NewRound creates new round, receives
func NewRound(triumph string) *round {
	return &round{
		triumph: triumph,
	}
}

// GetWinner returns the round winner
func (r *round) GetWinner() (p *cardPlayed) {
	if r.pos != 3 {
		// TODO QUITAR
		panic(1)
	}
	if err := r.checkWinner(); err != nil {
		panic(1)
	}
	return r.winner
}

func (r *round) playedCard(p *state.Player, c *state.Card) {
	pc := &cardPlayed{
		Player: p,
		Card:   c,
	}

	if r.pos == 0 {

		r.suit = c.Suit
	}
	r.played[r.pos] = pc
	r.cardsPlayed[r.pos] = c
	// Sum of points
	r.points += c.Points
	if r.pos < 3 {
		r.pos++
	}
}

// Returns true if there is a winner
func (r *round) checkWinner() error {

	if r.pos < 3 {
		return fmt.Errorf("not enough cards played")
	}
	winnerCard := r.played[0]
	for _, c := range r.played {

		// If they same suits and not wins
		if winnerCard.SameSuit(c.Card) && !winnerCard.Wins(c.Card) {
			winnerCard = c
		} else if c.IsTriumph(r.triumph) && !winnerCard.SameSuit(c.Card) {
			winnerCard = c
		}
	}
	r.winner = winnerCard
	return nil
}

func (r *round) Points() int {
	return r.points
}
