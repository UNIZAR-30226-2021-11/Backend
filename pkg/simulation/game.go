package simulation

import (
	"Backend/pkg/state"
)

const (
	starting = iota
	// Turnos de los jugadores
	cardsDealed

	t1
	t2
	t3
	t4
	// Check sings
	singing
	// Check triumph changes
	swap7

	// Check if there is a winner
	checkWinner
	checkWinnerVueltas
	// Has a winner
	winner
)

// Game has 10 rounds
type Game struct {
	players      *state.Ring
	rounds       [10]*round
	triumph      string
	currentRound int
	GameState    GameState `json:"game_state"`

	team1Points int
	team2Points int
	deck        *state.Deck
	cardsPlayed int
	// notifies the winner
	hasWinner chan<- struct{}

	// Input
	stop    <-chan struct{}
	newCard <-chan Event
}

type GameState struct {
	PointsTeamA int `json:"points_team_a"`
	PointsTeamB int `json:"points_team_b"`
	PointsSingA int `json:"points_sing_a"`
	PointsSingB int `json:"points_sing_b"`

	currentState int
	CurrentRound int `json:"current_round"`

	Vueltas bool `json:"vueltas"`

	Players *state.Ring `json:"players"`
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

func (g *Game) roundHasWinner() bool {
	return g.rounds[g.currentRound].checkWinner()
}

// Process a new card played
func (g *Game) processCard(c *state.Card) {

	switch g.GameState.currentState {
	case t1:
		g.cardPlayed(c)
		g.GameState.currentState = t2
	case t2:
		g.cardPlayed(c)

		g.GameState.currentState = t3
	case t3:
		g.cardPlayed(c)

		g.GameState.currentState = t4
	case t4:
		g.cardPlayed(c)

		if g.GameState.Vueltas {
			g.GameState.currentState = checkWinnerVueltas
			g.checkWinnerVueltas()
		} else {
			g.GameState.currentState = singing
			g.singingState()
		}
	}
	// Cartas jugadas
	// Nueva ronda
	// Esperar Cantes
	// Cambiar 7
	// Repartir cartas
}
func (g *Game) checkWinnerVueltas() {

}
func (g *Game) singingState() {

}
func (g *Game) processSing(playerId int, suit string) {

}

func (g *Game) changeCard(playerId int) {

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
		switch ev := c.(type) {
		case CardPlayedEvent:
			g.processCard(state.CreateCard(ev.suit, ev.val))

		case SingEvent:
			g.processSing(ev.playerId, ev.singSuit)

		case CardChangeEvent:
			g.changeCard(ev.playerId)
		}

		// Process new state
		// Repartir cartas
		// Actualizar cantes
		// Actualizar cambiar el 7
	}
}

// Handlers

func (g *Game) HandleCardPlayed() {

}

func (g *Game) HandleSing() {

}

func (g *Game) HandleNoSing() {

}

func (g *Game) HandleChangedCard() {

}

func (g *Game) GetState() {

}
