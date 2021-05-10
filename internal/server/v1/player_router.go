package v1

import (
	"Backend/internal/middleware"
	"Backend/pkg/player"
	"Backend/pkg/response"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

type PlayerRouter struct {
	Repository player.Repository
}

// CreateHandler Create a new player.
func (pr *PlayerRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var p player.Player
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = pr.Repository.Create(ctx, &p)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusCreated, response.Map{"player": p})
}

// Routes returns player router with each endpoint.
func (pr *PlayerRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Authorizator)


	r.Post("/", pr.CreateHandler)

	return r
}