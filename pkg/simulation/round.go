package simulation

import (
	"Backend/pkg/state"
	"fmt"
)

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
		return nil
		//panic(1)
	}
	if err := r.checkWinner(); err != nil {
		panic(1)
	}
	return r.winner
}

func (r *round) GetCardsPlayed() []*state.Card {

	var played []*state.Card
	for _, c := range r.played {
		if c != nil {
			played = append(played, c.Card)
		}
	}

	return played
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
	r.points += state.GetPoints(c.Val)
	if r.pos < 3 {
		r.pos++
	}
	r.updateWinnerCard()
}

// Returns true if there is a winner
func (r *round) checkWinner() error {

	if r.pos < 3 {
		return fmt.Errorf("not enough cards played")
	}
	r.updateWinnerCard()
	return nil
}

func (r *round) updateWinnerCard() {
	winnerCard := r.played[0]
	for _, c := range r.played {
		if c == nil {
			continue
		}
		// If they same suits and not wins
		if winnerCard.SameSuit(c.Card) && !winnerCard.Wins(c.Card) {
			winnerCard = c
		} else if c.IsTriumph(r.triumph) && !winnerCard.SameSuit(c.Card) {
			winnerCard = c
		}
	}
	r.winner = winnerCard
}

// Points returns the points gained in this round.
func (r *round) Points() int {
	return r.points
}

func (r *round) LastCard() *state.Card {
	return r.played[r.pos].Card
}

// CanPlayCards checks whether this suit can be played at the current stage.
func (r *round) CanPlayCards(arrastre bool, cs []*state.Card) {
	if !arrastre || r.pos == 0 {

		setPlayable(true, cs)
		return
	}

	setPlayable(false, cs)

	sameSuitCards := r.GetSuitCards(r.suit, cs)
	triumphCards := r.GetSuitCards(r.triumph, cs)

	lastWinner := r.winner.Card
	// If has cards of the same suit
	if len(sameSuitCards) > 0 {
		canWin := false
		for _, c := range sameSuitCards {
			if c.Wins(lastWinner) {
				canWin = true
				c.AllowPlay(true)
			}
		}
		if !canWin {
			setPlayable(true, sameSuitCards)
		}
		return
	}

	// If has triumph and can kill
	if len(triumphCards) > 0 {
		canWin := false
		killTriumph := lastWinner.IsTriumph(r.triumph)
		if killTriumph {
			for _, c := range triumphCards {
				// Last winner card was a triumph, has to win it
				if lastWinner.IsTriumph(r.triumph) {
					if c.Wins(lastWinner) {
						canWin = true
						c.AllowPlay(true)
					}
				}
			}
		} else {
			// Can play all triumphs
			canWin = true
			setPlayable(true, triumphCards)
		}
		// If the player has a card that can kill
		if canWin {
			return
		}
	}
	// Allow all cards
	setPlayable(true, cs)

}

func setPlayable(playable bool, cards []*state.Card) {
	for _, c := range cards {
		if c != nil {
			c.AllowPlay(playable)
		}
	}
}

func (r *round) GetSuitCards(suit string, cs []*state.Card) (sameSuit []*state.Card) {

	for _, c := range cs {
		if c.Suit == suit {
			sameSuit = append(sameSuit, c)
		}
	}

	return sameSuit
}
