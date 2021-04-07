package v1

import (
	"net/http"

	"Backend/internal/data"
	"github.com/go-chi/chi"
)

// New returns the API V1 Handler with configuration.
func New() http.Handler {
	r := chi.NewRouter()

	ur := &UserRouter{
		Repository: &data.UserRepository{
			Data: data.New(),
		},
	}

	r.Mount("/users", ur.Routes())

	gr := &GameRouter{
		Repository: &data.GameRepository{
			Data: data.New(),
		},
	}

	r.Mount("/games", gr.Routes())

	return r
}
