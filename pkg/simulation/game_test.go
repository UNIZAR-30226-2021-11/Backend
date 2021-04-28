package simulation

import (
	"Backend/pkg/state"
	"encoding/json"
	"testing"
)

var (
	ID_NEW_PLAYER uint32
	PAIR          uint32
	USERNAME      = "PEPE"
)

func TestCantar(t *testing.T) {
	p := state.CreatePlayer(0, 0, "pepe")
	p.DealCards([6]*state.Card{
		state.CreateCard(state.SUIT1, 10),
		state.CreateCard(state.SUIT2, 10),
		state.CreateCard(state.SUIT3, 10),
		state.CreateCard(state.SUIT1, 12),
		state.CreateCard(state.SUIT2, 12),
		state.CreateCard(state.SUIT3, 12),
	})
	var s Sings
	s.initialize()
	suits, canSing := p.HasSing()
	t.Run("player must sing", func(t *testing.T) {

		if !canSing {
			t.Errorf("got %v, want %v", canSing, !canSing)
		}
	})

	t.Run("sing 1 must be allowed", func(t *testing.T) {
		suit, canSing := s.canSign(suits)
		if !canSing {
			t.Errorf("got %v, want %v", canSing, !canSing)
		}

		if suit != state.SUIT1 {
			t.Errorf("got %v, want %v", suit, state.SUIT1)
		}
	})
	s.singedSuit(state.SUIT1)
	t.Run("sing 2 must be allowed", func(t *testing.T) {
		suit, canSing := s.canSign(suits)
		if !canSing {
			t.Errorf("got %v, want %v", canSing, !canSing)
		}

		if suit != state.SUIT2 {
			t.Errorf("got %v, want %v", suit, state.SUIT2)
		}
	})
	s.singedSuit(state.SUIT2)
	t.Run("sing 3 must be allowed", func(t *testing.T) {
		suit, canSing := s.canSign(suits)
		if !canSing {
			t.Errorf("got %v, want %v", canSing, !canSing)
		}

		if suit != state.SUIT3 {
			t.Errorf("got %v, want %v", suit, state.SUIT3)
		}
	})

	s.singedSuit(state.SUIT3)

	t.Run("no sings must be allowed", func(t *testing.T) {
		_, canSing := s.canSign(suits)
		if canSing {
			t.Errorf("got %v, want %v", canSing, !canSing)
		}
	})
}

func TestInitGame(t *testing.T) {

	var g *Game
	ps := createTestPlayers()
	t.Run("new game", func(t *testing.T) {
		g = NewGame(ps)
	})

	t.Run("check only one can play", func(t *testing.T) {
		count := 0
		for _, p := range g.GameState.Players.All {
			if p.CanPlay {
				count++
			}
		}
		if count != 1 {
			t.Errorf("got %v, want %v", count, 1)
		}
	})
	g.GameState.TriumphCard = g.deck.PickCard()
	t.Run("serialization json", func(t *testing.T) {
		b, err := json.Marshal(g.GameState)
		if err != nil {
			t.Error("error marshaling")
		}
		t.Log(string(b))

	})
}

func TestPlayOneCard(t *testing.T) {
	ps := createTestPlayers()
	g := NewGame(ps)
	initial := g.GameState.Players.Current()
	c := initial.Cards[0]

	t.Run("play one card", func(t *testing.T) {

		g.HandleCardPlayed(c)

		got := g.rounds[g.currentRound].cardsPlayed[0]
		if !got.Equals(c) {
			t.Errorf("got %v, want %v", got, c)
		}
	})

	want := c.Suit
	t.Run("correct suit", func(t *testing.T) {
		if g.rounds[g.currentRound].suit != want {
			t.Errorf("got %v, want %v",
				g.rounds[g.currentRound].triumph,
				want)
		}
	})

	t.Run("updated current player", func(t *testing.T) {
		current := g.GameState.Players.Current()
		if initial.Id == current.Id {
			t.Errorf("got same id")
		}
	})
}

func TestPlayRound(t *testing.T) {
	ps := createTestPlayers()
	g := NewGame(ps)

	var cardsPlayed []*state.Card
	r := NewRound(g.deck.GetTriumph())

	for i := 0; i < 4; i++ {
		c := g.GameState.Players.Current().PickRandomCard(0)
		g.HandleCardPlayed(c)
		r.playedCard(g.GameState.Players.Current(), c)
		cardsPlayed = append(cardsPlayed, c)
	}
	r.checkWinner()
	t.Run("check same winner", func(t *testing.T) {
		if !r.winner.Equals(g.rounds[0].winner.Card) {
			t.Errorf("got %v, want %v", g.rounds[0].winner, r.winner)
		}
	})

	t.Run("same points", func(t *testing.T) {
		want := r.Points()
		got := g.rounds[g.currentRound-1].Points()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("winner team has points", func(t *testing.T) {
		if g.winnerLastRound.Pair == 0 {
			t.Errorf("not updated winner pair")
		}
		want := r.Points()
		var got int
		if g.winnerLastRound.Pair == TeamA {
			got = g.GameState.PointsTeamA
		} else {
			got = g.GameState.PointsTeamB
		}
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("check integrity", func(t *testing.T) {
		if g.winnerLastRound.Player != g.GameState.Players.Current() {
			t.Errorf("Distinto jugador")
		}
		if g.currentRound != 1 {
			t.Errorf("mala ronda")
		}
		if g.GameState.Players.Current() != g.winnerLastRound.Player {
			t.Errorf("Distinto jugador")
		}
	})

}

func TestPlay2Round(t *testing.T) {
	ps := createTestPlayers()
	g := NewGame(ps)

	var cardsPlayed []*state.Card
	r := NewRound(g.deck.GetTriumph())

	for i := 0; i < 4; i++ {
		p := g.GameState.Players.Current()
		c := p.PickRandomCard(0)
		g.HandleCardPlayed(c)
		r.playedCard(p, c)
		cardsPlayed = append(cardsPlayed, c)
	}
	r.checkWinner()
	t.Run("check same winner", func(t *testing.T) {
		if !r.winner.Equals(g.rounds[0].winner.Card) {
			t.Errorf("got %v, want %v", g.rounds[0].winner, r.winner)
		}
	})

	t.Run("next round has correct player", func(t *testing.T) {
		if g.GameState.Players.Current().Id != r.winner.Player.Id {
			t.Errorf("got %v, want %v",
				g.GameState.Players.Current().Id, r.winner.Player.Id)
		}
	})

}

func TestDealedCards(t *testing.T) {
	ps := createTestPlayers()
	g := NewGame(ps)

	var cardsPlayed []*state.Card
	r := NewRound(g.deck.GetTriumph())

	for i := 0; i < 4; i++ {
		c := g.GameState.Players.Current().PickRandomCard(0)
		g.HandleCardPlayed(c)
		r.playedCard(g.GameState.Players.Current(), c)
		cardsPlayed = append(cardsPlayed, c)
	}
	if checkDistinctCards(t, g) {
		t.Errorf("players have same cards")
	}

}

func checkDistinctCards(t *testing.T, g *Game) bool {
	for _, p1 := range g.GameState.Players.All {
		for _, c1 := range p1.Cards {
			for _, p2 := range g.GameState.Players.All {
				if p1.Id != p2.Id {
					for _, c2 := range p2.Cards {
						if c1 != nil && c2 != nil && c1.Equals(c2) {
							t.Errorf("%v, %v", c1, c2)
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func TestPlayAllRounds(t *testing.T) {
	ps := createTestPlayers()
	g := NewGame(ps)

	var cardsPlayed []*state.Card

	for rs := 0; rs < 10; rs++ {
		for i := 0; i < 4; i++ {
			p := g.GameState.Players.Current()
			c := p.PickCard(0)
			g.HandleCardPlayed(c)
			cardsPlayed = append(cardsPlayed, c)
		}

		t.Run("dealed distinct cards", func(t *testing.T) {
			checkDistinctCards(t, g)
		})

		t.Run("just one player can play", func(t *testing.T) {
			count := 0
			for _, p := range ps {
				if p.CanPlay {
					count++
				}
			}
			if count > 1 {
				t.Errorf("more than 1 player can play")
			}
		})
	}

	t.Run("game state ended", func(t *testing.T) {
		if g.winnerPair != 0 && g.GameState.currentState != ended {
			t.Errorf("got %v, want %v", g.GameState.currentState, ended)
		}
		if g.winnerPair != 0 {
			t.Logf("Winner team: %v", g.winnerPair)
		}
		t.Logf("Puntos A: %v Puntos B: %v Sum: %v",
			g.GetTeamPoints(1),
			g.GetTeamPoints(2),
			g.GetTeamPoints(1)+g.GetTeamPoints(2))
	})

	t.Run("correct cards played", func(t *testing.T) {
		if len(cardsPlayed) != 40 {
			t.Errorf("got %v, want 40", len(cardsPlayed))
		}
	})

	t.Run("correct total points", func(t *testing.T) {
		if sumPoints(g.rounds) != 120 {
			t.Errorf("got %v, want %v", sumPoints(g.rounds), 120)
		}

		sum := 0
		for _, c := range cardsPlayed {
			sum += c.Points
		}
		if sum != 120 {
			t.Errorf("got %v, want %v", sum, 120)
		}
	})

	t.Run("needs rematch", func(t *testing.T) {
		if g.GetTeamPoints(TeamA) < 100 && g.GetTeamPoints(TeamB) < 100 {
			t.Logf("needs rematch")
			if !g.GameState.Vueltas {
				t.Errorf("rematch want %v, got %v", true, g.GameState.Vueltas)
			}
		} else {
			t.Logf("doesn't need rematch")
		}
	})
}

func hasBeenPlayed(played []*state.Card, c *state.Card) bool {
	for _, c2 := range played {
		if c2.Equals(c) {
			return true
		}
	}
	return false
}

func TestInitialCardDealing(t *testing.T) {
	g := NewGame(createTestPlayers())

	g.initialCardDealing()
	b, err := json.Marshal(g.GameState)
	if err != nil {
		t.Error("error marshaling")
	}
	t.Log(string(b))

	t.Run("correct number of cards", func(t *testing.T) {

		for player := 0; player < 4; player++ {
			p := g.GameState.Players.GetN(player).Cards
			count := 0
			for i := 0; i < 6; i++ {
				if p[i] != nil {
					count++
				}
			}
			if count != 6 {
				t.Errorf("got %v, want %v", count, 6)
			}
		}
	})
}

func createTestPlayers() []*state.Player {
	var players []*state.Player
	for i := 0; i < 4; i++ {
		players = append(players, CreateTestPlayer())
	}
	return players
}

func CreateTestPlayer() *state.Player {

	defer func() {
		ID_NEW_PLAYER++
		PAIR++
	}()

	return state.CreatePlayer(ID_NEW_PLAYER, PAIR%2+1, USERNAME)
}

func sumPoints(rs [10]*round) (sum int) {
	for _, r := range rs {
		sum += r.points
	}
	return sum
}
