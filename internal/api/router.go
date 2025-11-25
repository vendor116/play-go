package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/vendor116/playgo/internal/generated"
)

func GetRouter(server generated.ServerInterface) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/v1/info", server.GetInfo)

	return r
}
