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
	futureGames     map[uint32][]*state.Player
	pausedGames     map[uint32]*simulation.Game
	games           map[uint32]*simulation.Game
}

type GameData struct {
	Status string               `json:"status,omitempty"`
	Game   simulation.GameState `json:"game_state,omitempty"`
}

func NewSimulationRepository(eventDispatcher *events.EventDispatcher) *SimulationRepository {
	return &SimulationRepository{
		eventDispatcher: eventDispatcher,
		futureGames:     make(map[uint32][]*state.Player),
		games:           make(map[uint32]*simulation.Game),
		pausedGames:     make(map[uint32]*simulation.Game),
	}
}

func (sr *SimulationRepository) HandleGameCreate(gameCreateEvent *events.GameCreate) {
	log.Printf("User %d trying to create game %d\n", gameCreateEvent.PlayerID, gameCreateEvent.GameID)

	gameId := gameCreateEvent.GameID
	var players []*state.Player
	sr.futureGames[gameId] = players

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
	players, ok := sr.futureGames[gameId]
	if !ok {
		log.Printf("Game %d not found\n", gameId)
		return
	}

	newPlayers := append(players, player)

	if len(newPlayers) == 4 {
		sr.startNewGame(newPlayers, gameId)
	}
}

func (sr *SimulationRepository) startNewGame(players []*state.Player, gameId uint32) {
	firstPair := players[0].Pair

	for _, player := range players {
		if player.Pair != firstPair {
			player.Pair = 2
		} else {
			player.Pair = 1
		}
	}

	game := simulation.NewGame(players)

	sr.games[gameId] = game
	//delete(sr.futureGames, gameId)

	log.Printf("Game %v: Triumph is %v", gameId, game.GameState.TriumphCard.Suit)

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

func (sr *SimulationRepository) HandleVotePause(votePauseEvent *events.VotePause) {
	game, ok := sr.games[votePauseEvent.GameID]
	if !ok {
		log.Printf("Game %d not found\n", votePauseEvent.GameID)
		return
	}

	if votePauseEvent.Vote {
		sr.sendNewState(game.GameState, STATUS_PAUSED, game.GetPlayersID())
	}

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

	log.Printf("Client %v Game %v: Played card: %v", cardPlayedEvent.PlayerID, cardPlayedEvent.GameID, cardPlayedEvent.Card)

	sr.sendNewState(game.GameState, STATUS_NORMAL, game.GetPlayersID())
}

func (sr *SimulationRepository) HandleCardChanged(cardChangedEvent *events.CardChanged) {
	game, ok := sr.games[cardChangedEvent.GameID]
	if !ok {
		log.Printf("Game %d not found\n", cardChangedEvent.GameID)
		return
	}

	game.HandleChangedCard(cardChangedEvent.Changed)

	log.Printf("Client %v Game %v: Changed card: %v", cardChangedEvent.PlayerID,
		cardChangedEvent.GameID, cardChangedEvent.Changed)

	sr.sendNewState(game.GameState, STATUS_NORMAL, game.GetPlayersID())
}

func (sr *SimulationRepository) HandleSing(singEvent *events.Sing) {
	game, ok := sr.games[singEvent.GameID]
	if !ok {
		log.Printf("Game %d not found\n", singEvent.GameID)
		return
	}

	game.HandleSing(singEvent.Suit, singEvent.HasSinged)

	log.Printf("Client %v Game %v: Changed card: %v %v", singEvent.PlayerID,
		singEvent.GameID, singEvent.Suit, singEvent.HasSinged)

	sr.sendNewState(game.GameState, STATUS_NORMAL, game.GetPlayersID())
}

func (sr *SimulationRepository) sendNewState(game simulation.GameState,
	status string, clients []uint32) {
	gameData := &GameData{
		Status: status,
		Game:   game,
	}

	event := &events.StateChanged{
		ClientsID: clients,
		GameData:  gameData,
	}
	sr.eventDispatcher.FireStateChanged(event)
}
