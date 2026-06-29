package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HealthCheckRoutes(r chi.Router) {
	r.Get("/healthcheck", func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
