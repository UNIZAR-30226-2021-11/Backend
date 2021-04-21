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
		GameID:   gameCreateEvent.GameID,
		UserName: "usuario de prueba",
	})
}

func (sr *SimulationRepository) HandleUserJoined(userJoinedEvent *events.UserJoined) {
	gameId := userJoinedEvent.GameID
	player := &state.Player{
		Id: userJoinedEvent.PlayerID,
	}

	sr.futureGames[gameId] <- player
	players := sr.futureGames[gameId]

	if len(players) == 4 {
		var ps []*state.Player
		for i := 0; i < 4; i++ {
			ps = append(ps, <-players)
		}
		game := simulation.NewGame(ps)
		//game.InitGame()
		sr.games[gameId] = game
		event := &events.StateChanged{
			ClientsID: []uint32{1},
			Game:      game.GameState,
		}
		sr.eventDispatcher.FireStateChanged(event)
		delete(sr.futureGames, gameId)
	}
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
