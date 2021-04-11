package v1

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


// Routes returns game router with each endpoint.
func (gr *GameRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Authorizator)

	r.Get("/user/{userId}", gr.GetByUserHandler)

	r.Get("/", gr.GetAllHandler)

	r.Get("/{gameId}", gr.GetOneHandler)

	r.Post("/{userId}", gr.CreateHandler)

	//
	//r.Get("/{name}", gr.GetByName)
	//
	//r.Put("/{id}", gr.UpdateHandler)
	//
	//r.Delete("/{id}", gr.DeleteHandler)

	return r
}