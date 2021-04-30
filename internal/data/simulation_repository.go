package data

import (
	"Backend/pkg/events"
	"Backend/pkg/simulation"
	"Backend/pkg/state"
	"log"
)

const (
	STATUS_NORMAL = "normal"
	STATUS_VOTE   = "vote"
	STATUS_PAUSED = "paused"
)

type SimulationRepository struct {
	eventDispatcher *events.EventDispatcher
	futureGames     map[uint32]chan *state.Player
	pausedGames     map[uint32]*simulation.Game
	games           map[uint32]*simulation.Game
}

type GameData struct {
	status string               `json:"status,omitempty"`
	game   simulation.GameState `json:"game_state,omitempty"`
}

func NewSimulationRepository(eventDispatcher *events.EventDispatcher) *SimulationRepository {
	return &SimulationRepository{
		eventDispatcher: eventDispatcher,
		futureGames:     make(map[uint32]chan *state.Player),
		games:           make(map[uint32]*simulation.Game),
		pausedGames:     make(map[uint32]*simulation.Game),
	}
}

func (sr *SimulationRepository) HandleGameCreate(gameCreateEvent *events.GameCreate) {
	log.Printf("User %d trying to create game %d\n", gameCreateEvent.PlayerID, gameCreateEvent.GameID)

	gameId := gameCreateEvent.GameID
	sr.futureGames[gameId] = make(chan *state.Player, 4)

	sr.eventDispatcher.FireUserJoined(&events.UserJoined{
		PlayerID: gameCreateEvent.PlayerID,
		PairID:   gameCreateEvent.PairID,
		GameID:   gameCreateEvent.GameID,
		UserName: gameCreateEvent.UserName,
	})
}
func (sr *SimulationRepository) HandleUserJoined(userJoinedEvent *events.UserJoined) {
	gameId := userJoinedEvent.GameID
	player := &state.Player{
		Id:   userJoinedEvent.PlayerID,
		Pair: userJoinedEvent.PairID,
	}

	sr.futureGames[gameId] <- player
	players := sr.futureGames[gameId]

	if len(players) == 4 {
		sr.startNewGame(players, gameId)
	}
}

func (sr *SimulationRepository) startNewGame(playersChan chan *state.Player, gameId uint32) {
	var playersArray []*state.Player
	player := <-playersChan
	firstPair := player.Pair
	player.Pair = 1
	playersArray = append(playersArray, player)
	for i := 0; i < 3; i++ {
		player := <-playersChan
		if player.Pair != firstPair {
			player.Pair = 2
		} else {
			player.Pair = 1
		}
		playersArray = append(playersArray, player)
	}
	game := simulation.NewGame(playersArray)

	sr.games[gameId] = game
	delete(sr.futureGames, gameId)

	sr.sendNewState(game.GameState, STATUS_NORMAL, game.GetPlayersID())
}

func (sr *SimulationRepository) HandleGamePause(gamePauseEvent *events.GamePause) {
	game, ok := sr.games[gamePauseEvent.GameID]
	if !ok {
		log.Printf("Game %d not found\n", gamePauseEvent.GameID)
		return
	}

	opponents := game.GetOpponentsID(gamePauseEvent.PlayerID)
	sr.sendNewState(game.GameState, STATUS_VOTE, opponents)
}

func (sr *SimulationRepository) HandleUserLeft(userLeftEvent *events.UserLeft) {
	//TODO: change user by IA
}

func (sr *SimulationRepository) HandleCardPlayed(cardPlayedEvent *events.CardPlayed) {
	game, ok := sr.games[cardPlayedEvent.GameID]
	if !ok {
		log.Printf("Game %d not found\n", cardPlayedEvent.GameID)
		return
	}

	game.HandleCardPlayed(cardPlayedEvent.Card)

	sr.sendNewState(game.GameState, STATUS_NORMAL, game.GetPlayersID())
}

func (sr *SimulationRepository) HandleCardChanged(cardChangedEvent *events.CardChanged) {
	game, ok := sr.games[cardChangedEvent.GameID]
	if !ok {
		log.Printf("Game %d not found\n", cardChangedEvent.GameID)
		return
	}

	game.HandleChangedCard(cardChangedEvent.Changed)

	sr.sendNewState(game.GameState, STATUS_NORMAL, game.GetPlayersID())
}

func (sr *SimulationRepository) HandleSing(singEvent *events.Sing) {
	game, ok := sr.games[singEvent.GameID]
	if !ok {
		log.Printf("Game %d not found\n", singEvent.GameID)
		return
	}

	game.HandleSing(singEvent.Suit, singEvent.HasSinged)

	sr.sendNewState(game.GameState, STATUS_NORMAL, game.GetPlayersID())
}

func (sr *SimulationRepository) sendNewState(game simulation.GameState,
	status string, clients []uint32) {
	gameData := &GameData{
		status: status,
		game:   game,
	}

	event := &events.StateChanged{
		ClientsID: clients,
		GameData:  gameData,
	}
	sr.eventDispatcher.FireStateChanged(event)
}
