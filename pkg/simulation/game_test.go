package simulation

import (
	"Backend/pkg/state"
	"encoding/json"
	"testing"
)

var (
	ID_NEW_PLAYER uint32
	PAIR          = 0
)

func TestCantar(t *testing.T) {

}

func TestName(t *testing.T) {

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
	initial := g.currentPlayer
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
		current := g.currentPlayer
		if initial.Id == current.Id {
			t.Errorf("got same id")
		}
	})
}

func TestPlayRound(t *testing.T) {
	ps := createTestPlayers()
	g := NewGame(ps)

	var cardsPlayed []*state.Card
	r := NewRound(0, g.deck.GetTriumph())

	for i := 0; i < 4; i++ {
		c := g.currentPlayer.PickRandomCard(0)
		g.HandleCardPlayed(c)
		r.playedCard(c)
		cardsPlayed = append(cardsPlayed, c)
	}
	r.checkWinner()
	t.Run("check same winner", func(t *testing.T) {
		if !r.winner.Equals(g.rounds[0].winner) {
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
		if g.winnerPairLastRound == 0 {
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
}

func TestPlayAllRounds(t *testing.T) {
	ps := createTestPlayers()
	g := NewGame(ps)

	var cardsPlayed []*state.Card
	r := NewRound(0, g.deck.GetTriumph())

	for rounds := 0; rounds < 10; rounds++ {
		for i := 0; i < 4; i++ {
			c := g.currentPlayer.PickCard(0)
			g.HandleCardPlayed(c)
			r.playedCard(c)
			cardsPlayed = append(cardsPlayed, c)
		}
	}
	t.Run("game state ended", func(t *testing.T) {
		if g.winnerPair != 0 && g.GameState.currentState != ended {
			t.Errorf("got %v, want %v", g.GameState.currentState, ended)
		}
		if g.winnerPair != 0 {
			t.Logf("Winner team: %v", g.winnerPair)
		}
		t.Logf("Puntos A: %v Puntos B: %v Sum: %v",
			g.GetTeamPoints(1), g.GetTeamPoints(2), g.GetTeamPoints(1)+g.GetTeamPoints(2))
	})

	t.Run("correct total points", func(t *testing.T) {
		if sumPoints(g.rounds) != 120 {
			t.Errorf("got %v, want %v", sumPoints(g.rounds), 120)
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

	//t.Run("card played", func(t *testing.T) {
	//	g.HandleCardPlayed(0, state.CreateCard())
	//	g.rounds[g.currentRound].triumph
	//})
}

func createTestPlayers() []*state.Player {
	var players []*state.Player
	for i := 0; i < 4; i++ {
		players = append(players, CreateTestPlayer())
	}
	return players
}

func createGame(players []*state.Player) *Game {
	return NewGame(players)
}

func CreateTestPlayer() *state.Player {

	defer func() {
		ID_NEW_PLAYER++
		PAIR++
	}()

	return state.CreatePlayer(ID_NEW_PLAYER, PAIR%2+1)
}

func sumPoints(rs [10]*round) (sum int) {
	for _, r := range rs {
		sum += r.points
	}
	return sum
}
