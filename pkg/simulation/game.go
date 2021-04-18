package simulation

import (
	"Backend/pkg/state"
)

// Game has 10 rounds
type Game struct {
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

// Mejor un constructor, y una vez tenga el objeto llamo a startGame() p.e.
// El triunfo mejor genéralo aleatoriamente aquí, no tiene sentido pasarlo desde
//   el evento de creación de partida
func NewGame(chan *state.Player) *Game {
	return &Game{
		//TODO
	}

}

//Starts a new round
func (g *Game) newRound(firstPlayer *state.Player) {

	g.currentRound++
	g.rounds[g.currentRound] = NewRound(g.triumph)
	g.players.SetFirstPlayer(firstPlayer)

}
func (g *Game) cardPlayed(c *state.Card) {
	g.rounds[g.currentRound].playedCard(c)

}

// Process a new card played
func (g *Game) processCard(c *state.Card) {

	g.cardPlayed(c)
	// Cartas jugadas
	// Nueva ronda
	// Esperar Cantes
	// Cambiar 7
	// Repartir cartas
}

// StartGame starts a new Game with 10 rounds
func InitGame(p []*state.Player, triumph string) (g *Game) {

	g = &Game{
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

func (g *Game) Start() {

	for c := range g.newCard {
		g.processCard(c)
	}
}

func (g *Game) GetPlayersID() ([]uint32) {
	//TODO
	return []uint32{1,2,3,4}
}

func (g *Game) HandleCardPlayed(clientID uint32, card *state.Card) {
	//TODO
}