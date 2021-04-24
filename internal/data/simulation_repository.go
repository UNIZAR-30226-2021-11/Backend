package data

import (
	"Backend/pkg/events"
	"Backend/pkg/simulation"
	"Backend/pkg/state"
	"log"
)

type SimulationRepository struct {
	eventDispatcher *events.EventDispatcher
	futureGames     map[uint32]chan *state.Player
	games           map[uint32]*simulation.Game
}

func NewSimulationRepository(eventDispatcher *events.EventDispatcher) *SimulationRepository {
	return &SimulationRepository{
		eventDispatcher: eventDispatcher,
		futureGames:     make(map[uint32]chan *state.Player),
		games:           make(map[uint32]*simulation.Game),
	}
}

func (sr *SimulationRepository) HandleGameCreate(gameCreateEvent *events.GameCreate) {
	log.Printf("User %d trying to create game %d\n", gameCreateEvent.PlayerID, gameCreateEvent.GameID)

	gameId := gameCreateEvent.GameID
	sr.futureGames[gameId] = make(chan *state.Player, 4)

	sr.eventDispatcher.FireUserJoined(&events.UserJoined{
		PlayerID: gameCreateEvent.PlayerID,
		PairID: gameCreateEvent.PairID,
		GameID:   gameCreateEvent.GameID,
		UserName: gameCreateEvent.UserName,
	})
}

func (sr *SimulationRepository) HandleUserJoined(userJoinedEvent *events.UserJoined) {
	gameId := userJoinedEvent.GameID
	player := &state.Player{
		Id: userJoinedEvent.PlayerID,
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
	player :=  <-playersChan
	firstPair := player.Pair
	player.Pair = 1
	playersArray = append(playersArray, player)
	for i := 0; i < 3; i++ {
		player :=  <-playersChan
		if player.Pair != firstPair {
			player.Pair = 2
		} else {
			player.Pair = 1
		}
		playersArray = append(playersArray,player)
	}
	game := simulation.NewGame(playersArray)

	sr.games[gameId] = game
	event := &events.StateChanged{
		ClientsID: game.GetPlayersID(),
		Game:      game.GameState,
	}
	sr.eventDispatcher.FireStateChanged(event)
	delete(sr.futureGames, gameId)
}

func (sr *SimulationRepository) HandleUserLeft(userLeftEvent events.UserLeft) {
	//TODO
}

func (sr *SimulationRepository) HandleCardPlayed(cardPlayedEvent *events.CardPlayed) {
	gameId := cardPlayedEvent.GameID
	game := sr.games[gameId]

	game.HandleCardPlayed(cardPlayedEvent.Card)
	event := &events.StateChanged{
		ClientsID: game.GetPlayersID(),
		Game:      game,
	}
	sr.eventDispatcher.FireStateChanged(event)
}

func (sr *SimulationRepository) HandleCardChanged(cardChangedEvent *events.CardChanged) {
	//TODO
}

func (sr *SimulationRepository) HandleSing(singEvent *events.Sing) {
	//TODO
}
