package simulation

import "Backend/pkg/state"

// game has 10 rounds
type game struct {
	players      *state.Ring
	rounds       [10]*round
	triumph      string
	currentRound int

	team1Points int
	team2Points int
	deck        *state.Deck
	cardsPlayed int
	// notifies the winner
	hasWinner chan<- struct{}

	// Input
	stop    <-chan struct{}
	newCard <-chan *state.Card
}

//Starts a new round
func (g *game) newRound(firstPlayer *state.Player) {

	g.currentRound++
	g.rounds[g.currentRound] = NewRound(g.triumph)
	g.players.SetFirstPlayer(firstPlayer)

}
func (g *game) cardPlayed(c *state.Card) {
	g.rounds[g.currentRound].playedCard(c)

}

// Process a new card played
func (g *game) processCard(c *state.Card) {

	g.cardPlayed(c)
	// Cartas jugadas
	// Nueva ronda
	// Esperar Cantes
	// Cambiar 7
	// Repartir cartas
}

// StartGame starts a new game with 10 rounds
func InitGame(p []*state.Player, triumph string) (g *game) {

	g = &game{
		players:      state.NewPlayerRing(p),
		triumph:      triumph,
		currentRound: 0,
		deck:         state.NewDeck(),
	}
	g.deck.Shuffle()

	// Creates the first round
	g.rounds[0] = NewRound(triumph)
	// Set first player
	g.players.SetRandomFirstPlayer()
	return g
}

func (g *game) Start() {

	for c := range g.newCard {
		g.processCard(c)
	}
}
