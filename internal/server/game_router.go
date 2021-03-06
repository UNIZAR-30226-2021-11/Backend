package server

import (
	"Backend/internal/middleware"
	"Backend/pkg/game"
	"Backend/pkg/response"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

type GameRouter struct {
	Repository game.Repository
}

//CreateHandler Create and join a new game.
func (gr *GameRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var g game.Game
	err = json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = gr.Repository.Create(ctx, &g, uint(userID))
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), g.ID))
	response.JSON(w, r, http.StatusCreated, response.Map{"game": g})
}

//CreateHandler Create a new tournament game.
func (gr *GameRouter) CreateTournamentHandler(w http.ResponseWriter, r *http.Request) {
	var g game.Game
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = gr.Repository.CreateTournament(ctx, &g)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), g.ID))
	response.JSON(w, r, http.StatusCreated, response.Map{"game": g})
}

// GetAllHandler response all public started games.
func (gr *GameRouter) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	posts, err := gr.Repository.GetAll(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"games": posts})
}

// GetTournamentHandler response all tournament games.
func (gr *GameRouter) GetTournamentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	posts, err := gr.Repository.GetTournament(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"games": posts})
}

//GetByUserHandler response all ended games by user
func (gr *GameRouter) GetByUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	games, err := gr.Repository.GetByUser(ctx, uint(userID))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"games": games})
}

func (gr *GameRouter) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	gameIDStr := chi.URLParam(r, "gameId")

	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	game, err := gr.Repository.GetOne(ctx, uint(gameID))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"game": game})
}

func (gr *GameRouter) EndHandler(w http.ResponseWriter, r *http.Request) {
	var g game.Game
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = gr.Repository.End(ctx, g)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, nil)
}

// Routes returns game router with each endpoint.
func (gr *GameRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Put("/end", gr.EndHandler)

	r.With(middleware.Authorizator).Get("/user/{userId}", gr.GetByUserHandler)

	r.With(middleware.Authorizator).Get("/", gr.GetAllHandler)

	r.With(middleware.Authorizator).Get("/tournament", gr.GetTournamentHandler)

	r.With(middleware.Authorizator).Get("/{gameId}", gr.GetOneHandler)

	r.With(middleware.Authorizator).Post("/{userId}", gr.CreateHandler)

	r.With(middleware.Authorizator).Post("/tournament", gr.CreateTournamentHandler)

	return r
}
