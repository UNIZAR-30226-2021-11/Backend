package simulation

import (
	"Backend/pkg/state"
	"testing"
)

func TestRoundPlayed(t *testing.T) {
	r := NewRound(state.SUIT1)
	ps := createTestPlayers()

	cards := []*state.Card{
		state.CreateCard(state.SUIT4, 5),
		state.CreateCard(state.SUIT4, 1),
		state.CreateCard(state.SUIT2, 5),
		state.CreateCard(state.SUIT3, 3),
	}
	winnerCard := cards[1]
	t.Run("no_winner_at_start", func(t *testing.T) {

		err := r.checkWinner()
		if err == nil {
			t.Errorf("not enough cards played")
		}
		if r.winner != nil {
			t.Error("got a winner, shouldn't be")
		}
	})

	t.Run("cards_played_in_order", func(t *testing.T) {

		for i, card := range cards {
			if i != r.pos {
				t.Errorf("got %v, want %v", r.pos, i)
			}
			r.playedCard(ps[i], card)
		}
	})

	t.Run("check_correct_suit", func(t *testing.T) {

		suitWant := state.SUIT4
		if r.suit != suitWant {
			t.Errorf("got %v, want %v", r.suit, suitWant)
		}
	})

	t.Run("winner_at_end", func(t *testing.T) {

		err := r.checkWinner()
		if err != nil {
			t.Errorf("not enough cards played")
		}
		if r.winner == nil {
			t.Error("didn't get a winner, should be")
		}
	})

	t.Run("winner_correct_suit", func(t *testing.T) {

		if !r.winner.SameSuit(winnerCard) {
			t.Errorf("got %v, want %v", r.winner.Suit, winnerCard.Suit)
		}
	})

	t.Run("winner_correct_value", func(t *testing.T) {

		winnerValue := winnerCard.Points
		if r.winner.Points != winnerValue {
			t.Errorf("got %v, want %v", r.winner.Points, winnerValue)
		}
	})

	t.Run("round_has_correct_points", func(t *testing.T) {

		pointsWant := func() int {
			p := 0
			for _, c := range cards {
				p += c.Points
			}
			return p
		}()
		if r.Points() != pointsWant {

			t.Errorf("got %v, want %v", r.Points(), pointsWant)
		}
	})
}
func TestRoundPlayedWithTriumph(t *testing.T) {
	ps := createTestPlayers()

	r := NewRound(state.SUIT1)
	cards := []*state.Card{
		state.CreateCard(state.SUIT4, 5),
		state.CreateCard(state.SUIT4, 1),
		state.CreateCard(state.SUIT1, 4),
		state.CreateCard(state.SUIT1, 11),
	}
	winnerCard := cards[3]
	t.Run("no_winner_at_start", func(t *testing.T) {

		err := r.checkWinner()
		if err == nil {
			t.Errorf("not enough cards played")
		}
		if r.winner != nil {
			t.Error("got a winner, shouldn't be")
		}
	})

	t.Run("cards_played_in_order", func(t *testing.T) {

		for i, card := range cards {
			if i != r.pos {
				t.Errorf("got %v, want %v", r.pos, i)
			}
			r.playedCard(ps[i], card)
		}
	})

	t.Run("check_correct_suit", func(t *testing.T) {

		suitWant := state.SUIT4
		if r.suit != suitWant {
			t.Errorf("got %v, want %v", r.suit, suitWant)
		}
	})

	t.Run("winner_at_end", func(t *testing.T) {

		err := r.checkWinner()
		if err != nil {
			t.Errorf("not enough cards played")
		}
		if r.winner == nil {
			t.Error("didn't get a winner, should be")
		}
	})

	t.Run("winner_correct_suit", func(t *testing.T) {

		if !r.winner.SameSuit(winnerCard) {
			t.Errorf("got %v, want %v", r.winner.Suit, winnerCard.Suit)
		}
	})

	t.Run("winner_correct_value", func(t *testing.T) {

		winnerValue := winnerCard.Points
		if r.winner.Points != winnerValue {
			t.Errorf("got %v, want %v", r.winner.Points, winnerValue)
		}
	})

	t.Run("round_has_correct_points", func(t *testing.T) {

		pointsWant := func() int {
			p := 0
			for _, c := range cards {
				p += c.Points
			}
			return p
		}()
		if r.Points() != pointsWant {

			t.Errorf("got %v, want %v", r.Points(), pointsWant)
		}
	})
}
