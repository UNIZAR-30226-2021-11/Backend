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
	checkWinnerIdas
	checkWinnerVueltas
	// Has a winner
	ended
)

// Game has 10 rounds
type Game struct {
	rounds [10]*round
	//triumph      string
	currentRound int
	GameState    GameState `json:"game_state"`
	deck         *state.Deck

	pairCanSing     bool
	pairCanSwapCard bool

	currentPlayer   *state.Player
	winnerLastRound *state.Player
	topCard         *state.Card
}

type GameState struct {
	PointsTeamA int `json:"points_team_a"`
	PointsTeamB int `json:"points_team_b"`
	PointsSingA int `json:"points_sing_a"`
	PointsSingB int `json:"points_sing_b"`

	currentState  int
	CurrentRound  int    `json:"current_round"`
	CurrentPlayer uint32 `json:"current_player"`
	Vueltas       bool   `json:"vueltas"`

	Players *state.Ring `json:"players"`

	TriumphCard *state.Card `json:"triumph_card"`

	Arrastre bool `json:"arrastre"`
}

// NewGame returns a game in its initial state, with the deck shuffled
// and its first played picked
func NewGame(p []*state.Player) (g *Game) {
	r := state.NewPlayerRing(p)
	g = &Game{
		currentRound: 0,
		deck:         state.NewDeck(),
		GameState: GameState{
			PointsTeamA:  0,
			PointsTeamB:  0,
			PointsSingA:  0,
			PointsSingB:  0,
			currentState: 0,
			CurrentRound: 0,
			Vueltas:      false,
			Players:      r,
		},
	}
	g.deck.Shuffle()
	// Creates the first round
	g.rounds[0] = NewRound(g.deck.GetTriumph())
	// Set first player and deal initial cards
	g.GameState.Players.SetRandomFirstPlayer()
	g.initialCardDealing()

	first := g.GameState.Players.Current()
	g.GameState.CurrentPlayer = first.Id
	g.currentPlayer = first
	g.GameState.currentState = t1
	return g
}

// initial card dealing, 6 cards to each player
func (g *Game) initialCardDealing() {

	cards := g.deck.InitialPick()
	g.GameState.Players.InitialCardDealing(cards)
}

//Starts a new round
func (g *Game) newRound(firstPlayer *state.Player) {

	g.currentRound++
	g.rounds[g.currentRound] = NewRound(g.deck.GetTriumph())
	if g.currentRound > 6 {
		g.GameState.Arrastre = true
	}

	if !g.GameState.Arrastre {
		g.dealCards()
	}
	g.GameState.Players.SetFirstPlayer(firstPlayer)

}

// Process a card played, advances the player
func (g *Game) cardPlayed(c *state.Card) {
	g.rounds[g.currentRound].playedCard(c)
	g.GameState.Players.Current().PlayCard(c)
	g.GameState.Players.Next()
	g.currentPlayer = g.GameState.Players.Current()

}

// Checks for the round winner
func (g *Game) checkRoundWinner() {
	_, winnerPos := g.rounds[g.currentRound].checkWinner()
	g.winnerLastRound = g.GameState.Players.GetN(winnerPos)
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

		// TODO actualizar ganador de ronda
		g.checkRoundWinner()
		points := g.rounds[g.currentRound].Points()
		if g.winnerLastRound.Pair == 0 {
			g.GameState.PointsTeamA += points
		} else {
			g.GameState.PointsTeamB += points
		}

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

	// TODO comprobar si se ha ganado
	winner := false
	//
	if !winner {
		g.GameState.currentState = singing
		g.singingState()
	} else {
		g.GameState.currentState = ended
		g.ended()
	}
}
func (g *Game) singingState() {

	if !g.pairCanSing {
		g.GameState.currentState = swap7
		g.swapCard()
	} else {
		g.checkWinnerIdas()
	}

}
func (g *Game) processSing(suit string) {
	g.GameState.currentState = singing
	//if g.GameState.
}

func (g *Game) changeCard(hasChanged bool) {

	if hasChanged {
		seven := g.currentPlayer.GetSeven(g.GameState.TriumphCard.Suit)
		g.deck.ChangeCard(seven)
		g.currentPlayer.ChangeCard(g.topCard)

	}

}

func (g *Game) swapCard() {

	// TODO comprobar si la pareja ganadora puede cambiar
	if !g.pairCanSwapCard {
		if g.currentRound == 9 {

			g.GameState.currentState = checkWinnerIdas
			g.checkWinnerIdas()
		} else {
			g.GameState.currentState = t1
			g.newRound(g.winnerLastRound)
		}
	}
}

func (g *Game) checkWinnerIdas() {

	// TODO comprobar ganador idas
	winnerIdas := true
	if winnerIdas {
		g.GameState.currentState = ended
		g.ended()
	} else {
		//
		g.GameState.Vueltas = true
		g.restart()
	}
}

// Handlers

func (g *Game) HandleCardPlayed(card *state.Card) {
	g.processCard(card)
}

func (g *Game) HandleSing(suit string, hasSinged bool) {
	if hasSinged {
		g.processSing(suit)
	}
}

func (g *Game) HandleChangedCard(changedCard bool) {
	g.changeCard(changedCard)

}

// GetPlayersID returns the ids of all players
func (g *Game) GetPlayersID() []uint32 {
	return g.GameState.Players.GetPlayersIds()
}

func (g *Game) ended() {

}

// restart puts the game in the initial state, saves the current point
func (g *Game) restart() {

}

// Deals a card to each player
func (g *Game) dealCards() {
	cards := g.deck.Pick4Cards()
	g.GameState.Players.DealCards(cards)
}
