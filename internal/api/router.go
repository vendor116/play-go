package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/vendor116/playgo/internal/generated"
)

func GetRouter(server generated.ServerInterface) chi.Router {
	return chi.NewRouter().
		Route("/v1", func(r chi.Router) {
			r.Get("/info", server.GetInfo)
		})
}
