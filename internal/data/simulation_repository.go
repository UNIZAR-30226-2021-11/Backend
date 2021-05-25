package data

import (
	"Backend/pkg/ai"
	"Backend/pkg/events"
	"Backend/pkg/pair"
	"Backend/pkg/simulation"
	"Backend/pkg/state"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	ais             map[uint32]struct{}
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
		ais:             make(map[uint32]struct{}),
	}
}

func (sr *SimulationRepository) HandleSingleGameCreate(singleGameCreateEvent *events.SingleGameCreate) {
	log.Printf("User %d trying to create single game %d\n",
		singleGameCreateEvent.PlayerID, singleGameCreateEvent.GameID)

	gameId := singleGameCreateEvent.GameID
	var players []*state.Player
	sr.futureGames[gameId] = players

	var ais []*ai.Client
	for i := 1; i < 4; i++ {
		aiClient := ai.Create(uint32(i), uint32(i%2)+1, gameId)
		ais = append(ais, aiClient)
	}

	// Add the single user
	sr.eventDispatcher.FireUserJoined(&events.UserJoined{
		PlayerID: singleGameCreateEvent.PlayerID,
		PairID:   uint32(1),
		GameID:   singleGameCreateEvent.GameID,
		UserName: singleGameCreateEvent.UserName,
	})
	for _, c := range ais {
		c.Start()
		log.Printf("ai client: %d started", c.Id)
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
	log.Printf("User %d trying to join game %d\n", userJoinedEvent.PlayerID, userJoinedEvent.GameID)
	gameId := userJoinedEvent.GameID
	player := &state.Player{
		Id:         userJoinedEvent.PlayerID,
		InternPair: userJoinedEvent.PairID,
	}
	players, ok := sr.futureGames[gameId]
	if !ok {
		log.Printf("Game %d not found\n", gameId)
		return
	}

	players = append(players, player)
	sr.futureGames[gameId] = players

	if len(players) == 4 {
		sr.startGame(players, gameId)
	}
}

// startGame Starts a new game or restarts an existing game.
func (sr *SimulationRepository) startGame(players []*state.Player, gameId uint32) {
	pausedGame, isPaused := sr.pausedGames[gameId]

	if isPaused {
		sr.restartGame(pausedGame, gameId)
	} else {
		sr.startNewGame(players, gameId)
	}
}

func (sr *SimulationRepository) restartGame(game *simulation.Game, gameId uint32) {
	sr.games[gameId] = game

	delete(sr.pausedGames, gameId)
	delete(sr.futureGames, gameId)

	sr.sendNewState(gameId, game.GameState, STATUS_NORMAL, game.GetPlayersID())
}

func (sr *SimulationRepository) startNewGame(players []*state.Player, gameId uint32) {
	var newPlayers [4]*state.Player
	firstPair := players[0].InternPair
	counter1 := 0
	counter2 := 0

	for _, player := range players {
		if player.InternPair != firstPair {
			player.Pair = 2
			counter2++
			if counter2 == 1 {
				newPlayers[1] = player
			} else {
				newPlayers[3] = player
			}
		} else {
			player.Pair = 1
			counter1++
			if counter1 == 1 {
				newPlayers[0] = player
			} else {
				newPlayers[2] = player
			}
		}
	}

	game := simulation.NewGame(newPlayers[:])

	sr.games[gameId] = game
	delete(sr.futureGames, gameId)

	log.Printf("Game %v: Triumph is %v", gameId, game.GameState.TriumphCard.Suit)

	sr.sendNewState(gameId, game.GameState, STATUS_NORMAL, game.GetPlayersID())
}

func (sr *SimulationRepository) startSingleGame(player *state.Player, gameId uint32) {

	//game := simulation.NewGame(players)

}

func (sr *SimulationRepository) HandleGamePause(gamePauseEvent *events.GamePause) {
	game, ok := sr.games[gamePauseEvent.GameID]
	if !ok {
		log.Printf("Game %d not found\n", gamePauseEvent.GameID)
		return
	}

	opponents := game.GetOpponentsID(gamePauseEvent.PlayerID)
	sr.sendNewState(gamePauseEvent.GameID, game.GameState, STATUS_VOTE, opponents)
}

func (sr *SimulationRepository) HandleVotePause(votePauseEvent *events.VotePause) {
	gameId := votePauseEvent.GameID
	game, ok := sr.games[gameId]
	if !ok {
		log.Printf("Game %d not found\n", gameId)
		return
	}

	if votePauseEvent.Vote {
		// Let players know game is paused
		sr.sendNewState(votePauseEvent.GameID, game.GameState, STATUS_PAUSED, game.GetPlayersID())

		// Save the current game state
		sr.pausedGames[gameId] = game
		delete(sr.games, gameId)

		// For rejoining the game
		var players []*state.Player
		sr.futureGames[gameId] = players
	} else {
		sr.sendNewState(votePauseEvent.GameID, game.GameState, STATUS_NORMAL, game.GetPlayersID())
	}
}

func (sr *SimulationRepository) HandleUserLeft(userLeftEvent *events.UserLeft) {
	//TODO: change user by IA
	log.Printf("UserLeft %v", userLeftEvent.PlayerID)
	gid := userLeftEvent.GameID
	pid := userLeftEvent.PlayerID
	pairID := userLeftEvent.PairID
	game, ok := sr.games[gid]
	if !ok {
		log.Printf("Game %d not found\n", gid)
		return
	}
	_, aiExist := sr.ais[pid]
	if !aiExist {
		sr.ais[pid] = struct{}{}
		a := ai.Create(pid, pairID, gid)
		go a.TakeOver()
		sr.sendNewState(gid, game.GameState, STATUS_NORMAL, game.GetPlayersID())
	}
}

func (sr *SimulationRepository) HandleCardPlayed(cardPlayedEvent *events.CardPlayed) {
	game, ok := sr.games[cardPlayedEvent.GameID]
	if !ok {
		log.Printf("Game %d not found\n", cardPlayedEvent.GameID)
		return
	}

	game.HandleCardPlayed(cardPlayedEvent.Card)

	log.Printf("Client %v Game %v: Played card: %v", cardPlayedEvent.PlayerID, cardPlayedEvent.GameID, cardPlayedEvent.Card)

	sr.sendNewState(cardPlayedEvent.GameID, game.GameState, STATUS_NORMAL, game.GetPlayersID())
}

func (sr *SimulationRepository) HandleStateChanged(changed *events.StateChanged) {
	game, ok := sr.games[changed.GameID]
	if ok && game.GameState.Ended {
		pairId, points := game.GetWinningPair()
		sr.updatePairWon(pairId, true, points)

		delete(sr.games, changed.GameID)
		log.Printf("game %d ended", changed.GameID)
		for _, p := range changed.ClientsID {
			sr.eventDispatcher.FireUserLeft(&events.UserLeft{
				PlayerID: p,
				GameID:   changed.GameID,
				PairID:   0,
			})
			delete(sr.ais, p)
		}
	}
}

// updatePairWon updates the pair info in the API
func (sr *SimulationRepository) updatePairWon(pairID uint32, winned bool, gamePoints int) {
	pair := pair.Pair{
		Winned:     winned,
		GamePoints: gamePoints,
	}

	// initialize http client
	client := &http.Client{}

	// marshal User to json
	json, err := json.Marshal(pair)
	if err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	url := "http://localhost:" + port + "/api/v1/pairs/" + strconv.Itoa(int(pairID))
	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(json))
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)
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

	sr.sendNewState(cardChangedEvent.GameID, game.GameState, STATUS_NORMAL, game.GetPlayersID())
}

func (sr *SimulationRepository) HandleSing(singEvent *events.Sing) {
	game, ok := sr.games[singEvent.GameID]
	if !ok {
		log.Printf("Game %d not found\n", singEvent.GameID)
		return
	}

	game.HandleSing(singEvent.Suit, singEvent.HasSinged)

	log.Printf("Client %v Game %v: Singed suit: %v %v", singEvent.PlayerID,
		singEvent.GameID, singEvent.Suit, singEvent.HasSinged)

	sr.sendNewState(singEvent.GameID, game.GameState, STATUS_NORMAL, game.GetPlayersID())
}

func (sr *SimulationRepository) sendNewState(gameId uint32, game simulation.GameState,

	status string, clients []uint32) {
	gameData := &GameData{
		Status: status,
		Game:   game,
	}

	event := &events.StateChanged{
		ClientsID: clients,
		GameData:  gameData,
		GameID:    gameId,
	}
	sr.eventDispatcher.FireStateChanged(event)
}
