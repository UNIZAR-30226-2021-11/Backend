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

func TestRound_CanPlayCards(t *testing.T) {

	r := NewRound(state.SUIT1)
	/* Arrastre, Triunfo oros


	 */
	cards := []*state.Card{
		state.CreateCard(state.SUIT2, 3),
		state.CreateCard(state.SUIT4, 1),
		state.CreateCard(state.SUIT1, 4),
		state.CreateCard(state.SUIT1, 11),
	}
	arrastre := true

	cardsP2 := []*state.Card{
		state.CreateCard(state.SUIT4, 5),
		state.CreateCard(state.SUIT1, 2),
	}

	t.Run("no suit, triumph", func(t *testing.T) {
		r.playedCard(CreateTestPlayer(), cards[0])
		r.CanPlayCards(arrastre, cardsP2)
		if cardsP2[0].Playable {
			t.Errorf("shoulnd't be playable,got %v, want %v", cardsP2[0].Playable, false)
		}
	})

	r = NewRound(state.SUIT1)
	cardsP3 := []*state.Card{
		state.CreateCard(state.SUIT2, 5),
		state.CreateCard(state.SUIT1, 2),
	}

	t.Run("suit AND triumph", func(t *testing.T) {
		r.playedCard(CreateTestPlayer(), cards[0])
		r.CanPlayCards(arrastre, cardsP3)
		if !cardsP3[0].Playable && cardsP3[1].Playable {
			t.Errorf("just same suit card must be playable")
		}
	})

	r = NewRound(state.SUIT1)
	cardsP4 := []*state.Card{
		state.CreateCard(state.SUIT3, 5),
		state.CreateCard(state.SUIT4, 2),
	}

	t.Run("no suit AND no triumph", func(t *testing.T) {
		r.playedCard(CreateTestPlayer(), cards[0])
		r.CanPlayCards(arrastre, cardsP4)
		if !(cardsP4[0].Playable && cardsP4[1].Playable) {
			t.Errorf("should be all playable")
		}
	})

	r = NewRound(state.SUIT1)
	TRESCOPAS := state.CreateCard(state.SUIT2, 3)
	ASCOPAS := state.CreateCard(state.SUIT2, 1)
	cardsP5 := []*state.Card{
		state.CreateCard(state.SUIT2, 5),
		ASCOPAS,
	}

	t.Run("same suit, can win it", func(t *testing.T) {
		r.playedCard(CreateTestPlayer(), TRESCOPAS)
		r.CanPlayCards(arrastre, cardsP5)

		if cardsP5[0].Playable {
			t.Errorf("can play a not winner triumph")
		}
		if !cardsP5[1].Playable {
			t.Errorf("cannot play a winner triumph")
		}
	})

	r = NewRound(state.SUIT1)
	cardsP6 := []*state.Card{
		state.CreateCard(state.SUIT1, 1),
		state.CreateCard(state.SUIT1, 12),
	}

	t.Run("played triumph, can win it", func(t *testing.T) {
		// AS DE BASTOS
		r.playedCard(CreateTestPlayer(), state.CreateCard(state.SUIT4, 1))
		r.playedCard(CreateTestPlayer(), state.CreateCard(state.SUIT1, 3))
		r.CanPlayCards(arrastre, cardsP6)

		if !cardsP6[0].Playable {
			t.Errorf("winner card should be playable")
		}
		if cardsP6[1].Playable {
			t.Errorf("! winner triumph shouldn't be playable")
		}
		if t.Failed() {
			t.Logf("%v, %v", cardsP6[0], cardsP6[1])
		}
	})
}

func TestRound_GetCardsPlayed(t *testing.T) {
	r := NewRound(state.SUIT1)
	ps := createTestPlayers()

	cards := []*state.Card{
		state.CreateCard(state.SUIT4, 5),
		state.CreateCard(state.SUIT4, 1),
		state.CreateCard(state.SUIT2, 5),
		state.CreateCard(state.SUIT2, 1),
	}

	for i, card := range cards {
		r.playedCard(ps[i], card)
	}

	for i, card := range r.GetCardsPlayed() {
		if !card.Equals(cards[i]) {
			t.Errorf("not the same cards")
		}
	}
}
