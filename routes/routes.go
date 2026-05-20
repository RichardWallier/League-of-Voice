package routes

import (
	"fmt"
	"lov/handler"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(handlers *handler.Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	UsersRoutes(r, handlers)
	AuthRoutes(r, handlers)
	//UserMetricRoutes(r)

	fmt.Println("routers registered")
	return r
}
