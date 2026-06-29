package routes

import (
	"fmt"
	"lov/handler"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(handlers *handler.Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Apply a timeout middleware to all routes in this group
	r.Group(func(r chi.Router) {
		r.Use(middleware.Timeout(60 * time.Second))
		UsersRoutes(r, handlers.UserHandler)
		AuthRoutes(r, handlers.AuthHandler)
		HealthCheckRoutes(r)
	})

	SFURoutes(r, handlers.SFUHandler)
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
       http.ServeFile(w, req, "test-client.html")
   })	//UserMetricRoutes(r)

	fmt.Println("routers registered")
	return r
}
