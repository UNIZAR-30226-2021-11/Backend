package simulation

import (
	"Backend/pkg/state"
)

const (

	// Turnos de los jugadores
	t1 = iota
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

	TeamA = 1
	TeamB = 2
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

	winnerLastRound *cardPlayed
	topCard         *state.Card
	winnerLast10    uint32
	winnerPair      uint32

	// internal

	sings              Sings
	cardHasBeenChanged bool
}

type GameState struct {
	PointsTeamA int `json:"points_team_a"`
	PointsTeamB int `json:"points_team_b"`
	PointsSingA int `json:"points_sing_a"`
	PointsSingB int `json:"points_sing_b"`

	currentState int
	CurrentRound int `json:"current_round"`
	//CurrentPlayer uint32 `json:"current_player"`

	Players *state.Ring `json:"players"`

	TriumphCard *state.Card `json:"triumph_card"`

	Vueltas  bool `json:"vueltas"`
	Arrastre bool `json:"arrastre"`

	Ended      bool   `json:"ended"`
	WinnerPair uint32 `json:"winner_pair"`

	CardsPlayedRound []*state.Card `json:"cards_played_round"`
}

// Sings keeps track of singed suits
type Sings struct {
	// winner last round
	winnerPair int
	//
	sings map[string]bool
}

func (s *Sings) initialize() {
	s.winnerPair = 0
	s.sings = make(map[string]bool)
	s.sings[state.SUIT1] = false
	s.sings[state.SUIT2] = false
	s.sings[state.SUIT3] = false
	s.sings[state.SUIT4] = false
}

func (s *Sings) updateWinnerPair(wp int) {
	s.winnerPair = wp
}

func (s *Sings) singedSuit(suit string) {
	s.sings[suit] = true
}

// Checks whether this suits can be singed
func (s *Sings) canSign(suits []string) (string, bool) {
	for _, suit := range suits {
		canSing, ok := s.sings[suit]
		if ok && !canSing {
			return suit, true
		}
	}
	return "", false
}

// NewGame returns a game in its initial state, with the deck shuffled
// and its first played picked
func NewGame(p []*state.Player) (g *Game) {
	var s Sings
	s.initialize()
	r := state.NewPlayerRing(p)
	sings := make(map[string]bool)
	sings[state.SUIT1] = false
	sings[state.SUIT2] = false
	sings[state.SUIT3] = false
	sings[state.SUIT4] = false

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
		sings: s,
	}
	g.deck.Shuffle()

	// Set first player and deal initial cards
	g.GameState.Players.SetRandomFirstPlayer()
	g.GameState.TriumphCard = g.deck.GetTriumphCard()
	// Creates the first round
	g.rounds[0] = NewRound(g.deck.GetTriumph())
	g.initialCardDealing()

	for _, player := range g.GameState.Players.All {
		g.rounds[0].CanPlayCards(g.GameState.Arrastre, player.GetCards())
	}
	g.GameState.currentState = t1
	return g
}

// initial card dealing, 6 cards to each player
func (g *Game) initialCardDealing() {

	cards := g.deck.InitialPick()
	g.GameState.Players.InitialCardDealing(cards)
}

//Starts a new round
func (g *Game) newRound() {

	g.currentRound++
	g.GameState.CurrentRound++
	g.GameState.Players.SetFirstPlayer(g.winnerLastRound.Player)

	g.rounds[g.currentRound] = NewRound(
		g.deck.GetTriumph())

	if g.currentRound > 4 {
		g.GameState.Arrastre = true
	}

	if !g.GameState.Arrastre {
		g.dealCards()
	} else {
		//TODO COMPROBAR CARTAS
	}

}

// Process a new card played
func (g *Game) processCard(c *state.Card) {

	switch g.GameState.currentState {
	case t1:
		g.cardPlayed(c)
		g.GameState.CardsPlayedRound = g.rounds[g.currentRound].GetCardsPlayed()
		g.GameState.currentState = t2
	case t2:
		g.cardPlayed(c)
		g.GameState.CardsPlayedRound = g.rounds[g.currentRound].GetCardsPlayed()

		g.GameState.currentState = t3
	case t3:
		g.cardPlayed(c)
		g.GameState.CardsPlayedRound = g.rounds[g.currentRound].GetCardsPlayed()

		g.GameState.currentState = t4
	case t4:
		g.cardPlayed(c)
		g.GameState.CardsPlayedRound = g.rounds[g.currentRound].GetCardsPlayed()

		g.checkRoundWinner()
		g.updatePoints()
		g.updateSings()

		if !g.cardHasBeenChanged {
			g.updateChange()
		}

		if g.GameState.Vueltas {
			g.GameState.currentState = checkWinnerVueltas
			g.checkWinnerVueltas()
		} else {
			g.GameState.currentState = singing
			g.singingState()
		}
	}
}

// Process a card played, advances the player
func (g *Game) cardPlayed(c *state.Card) {
	current := g.GameState.Players.Current()
	r := g.rounds[g.currentRound]
	r.playedCard(current, c)
	current.PlayCard(c)

	for _, player := range g.GameState.Players.All {
		r.CanPlayCards(g.GameState.Arrastre, player.GetCards())
	}
	g.GameState.Players.Next()

}

// Checks for the round winner
func (g *Game) checkRoundWinner() {

	winner := g.rounds[g.currentRound].GetWinner()

	g.winnerLastRound = winner
	g.GameState.WinnerPair = winner.Pair
}

func (g *Game) updatePoints() {

	points := g.rounds[g.currentRound].Points()

	if g.winnerLastRound.Pair == TeamA {
		g.GameState.PointsTeamA += points
		if g.currentRound == 9 {
			g.GameState.PointsTeamA += 10
			g.winnerLast10 = TeamA
		}
	} else {
		g.GameState.PointsTeamB += points
		// 10 últimas
		if g.currentRound == 9 {
			g.GameState.PointsTeamB += 10
			g.winnerLast10 = TeamB

		}
	}
}

// Updates the singing state of a pair
func (g *Game) updateSings() {

	wp := g.winnerLastRound.Pair
	g.pairCanSing = false
	for _, p := range g.GameState.Players.All {
		p.CanSing = false
		if p.Pair == wp {
			suits, _ := p.HasSing()
			// If the player has a singing pair
			suit, allowed := g.sings.canSign(suits)

			if allowed {
				g.pairCanSing = true
				p.CanSing = true
				p.SingSuit = suit
				return
			}
		}
	}
}

func (g *Game) updateChange() {
	g.pairCanSwapCard = false
	for _, p := range g.GameState.Players.All {
		// Pair won round
		if p.Pair == g.winnerLastRound.Pair {
			seven := p.GetSeven(g.GameState.TriumphCard.Suit)
			if seven != nil {
				g.pairCanSwapCard = true
				p.CanChange = true
				break
			}
		}
	}
}

func (g *Game) checkWinnerVueltas() {

	// TODO comprobar si se ha ganado
	winner := g.checkWinner()
	//
	if !winner {
		g.GameState.currentState = singing
		g.singingState()
	} else {
		g.GameState.currentState = ended
		g.ended()
	}
}

func (g *Game) checkWinnerIdas() {

	winnerIdas := g.checkWinner()
	if winnerIdas {
		g.GameState.currentState = ended

		// Comprobar puntos de cada equipo
		g.ended()
	} else {
		//
		g.GameState.Vueltas = true
		g.restart()
	}
}

func (g *Game) singingState() {

	if !g.pairCanSing {
		g.GameState.currentState = swap7
		g.swapCard()

	}
}

func (g *Game) processSing(suit string) {
	g.GameState.currentState = singing
	g.sings.singedSuit(suit)

	if g.winnerLastRound.Pair == TeamA {
		if suit == g.GameState.TriumphCard.Suit {
			g.GameState.PointsSingA += 40
		} else {
			g.GameState.PointsSingA += 20
		}
	} else {
		if suit == g.GameState.TriumphCard.Suit {
			g.GameState.PointsSingB += 40
		} else {
			g.GameState.PointsSingB += 20
		}
	}
	g.updateSings()

	g.singingState()
}

func (g *Game) changeCard(hasChanged bool) {

	if hasChanged {
		g.cardHasBeenChanged = true
		triumph := g.GameState.TriumphCard.Suit
		for _, p := range g.GameState.Players.All {
			if p.Pair == g.winnerLastRound.Pair {
				seven := p.GetSeven(triumph)

				if seven != nil {
					last := g.deck.ChangeCard(seven)
					p.ChangeCard(triumph, last)
				}
			}
		}
	}
	for _, p := range g.GameState.Players.All {
		// Pair won round
		p.CanChange = false
	}
	g.pairCanSwapCard = false
	g.swapCard()
}

func (g *Game) swapCard() {

	if !g.pairCanSwapCard {
		if g.currentRound == 9 {

			g.GameState.currentState = checkWinnerIdas
			g.checkWinnerIdas()
		} else {
			g.GameState.currentState = t1
			g.newRound()
		}
	} else {
		g.GameState.currentState = swap7
	}
}

//
func (g *Game) checkWinner() bool {
	//SI una pareja no llega a 30 sin cantes, pierde
	if g.GameState.PointsTeamA < 30 {
		g.winnerPair = TeamB
		return true
	}
	if g.GameState.PointsTeamB < 30 {
		g.winnerPair = TeamA
		return true
	}
	// Si ambas superan 100, gana la que lleve 10 ultimas
	if g.GetTeamPoints(TeamA) > 100 && g.GetTeamPoints(TeamB) > 100 {
		g.winnerPair = g.winnerLast10
		return true
	}

	if g.GetTeamPoints(TeamA) > 100 {
		g.winnerPair = TeamA
		return true
	}
	if g.GetTeamPoints(TeamB) > 100 {
		g.winnerPair = TeamB
		return true
	}

	return false
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

// GetOpponentsID returns the ids of the other pair players
func (g *Game) GetOpponentsID(playerID uint32) []uint32 {
	return []uint32{1, 2}
}

// GetTeamPoints returns points for a team, even returns Team A, odd Team B
func (g *Game) GetTeamPoints(team int) (points int) {
	if team%2 == 0 {
		points = g.GameState.PointsTeamA + g.GameState.PointsSingA
	} else {
		points = g.GameState.PointsTeamB + g.GameState.PointsSingB
	}
	return points
}

func (g *Game) ended() {
	g.GameState.Ended = true
}

// restart puts the game in the initial state, saves the current point
func (g *Game) restart() {
	g.sings.initialize()
	g.deck = state.NewDeck()

	g.GameState.Vueltas = true
	g.GameState.Arrastre = false
	g.GameState.CurrentRound = 0
	g.currentRound = 0
	g.deck.Shuffle()

	for i := range g.rounds {
		g.rounds[i] = nil
	}

	g.pairCanSing = false
	g.pairCanSwapCard = false

	// Set first player and deal initial cards
	g.GameState.Players.SetFirstPlayer(g.winnerLastRound.Player)
	g.GameState.TriumphCard = g.deck.GetTriumphCard()

	// Creates the first round
	g.rounds[0] = NewRound(g.deck.GetTriumph())
	g.initialCardDealing()

	g.GameState.currentState = t1
}

// Deals a card to each player
func (g *Game) dealCards() {
	cards := g.deck.Pick4Cards()
	g.GameState.Players.DealCards(cards)
}
