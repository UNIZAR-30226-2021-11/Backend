package server

import (
	"Backend/internal/middleware"
	"Backend/pkg/pair"
	"Backend/pkg/response"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

// PairRouter is the router of the pairs.
type PairRouter struct {
	Repository pair.Repository
}

// UpdateWinnedHandler update a stored pair by id.
func (pr *PairRouter) UpdateWinnedHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var p pair.Pair
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = pr.Repository.UpdateWinned(ctx, uint(id))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, nil)
}

// Routes returns user router with each endpoint.
func (pr *PairRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.
		With(middleware.Authorizator).
		Put("/{id}", pr.UpdateWinnedHandler)

	return r
}
