package main

import (
	"context"
	"fmt"
	"lov/db"
	"lov/handler"
	"lov/repository"
	"lov/routes"
	"lov/service"
	"net/http"
)


func main() {
	fmt.Println("Server starting...")
	ctx := context.Background()

	db := db.NewPostgresDB(ctx)
	defer db.Cleanup()

	entities := repository.NewEntities(db)

	services := service.NewServices(entities)

	handlers := handler.SetupHandlers(services)

	routes := routes.SetupRoutes(handlers)

	fmt.Println("Server running on :3000...")
	if err := http.ListenAndServe(":3000", routes); err != nil {
		fmt.Println(err.Error())
	}
}
