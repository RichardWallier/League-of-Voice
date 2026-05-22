package main

import (
	"context"
	"log"
	"path/filepath"

	"lov/db"
	"lov/db/seeding"
)

func main() {
	ctx := context.Background()
	pg := db.NewPostgresDB(ctx)
	defer pg.Cleanup()

	files := []string{"roles.sql", "permissions.sql", "role_permissions.sql"}
	for _, f := range files {
		path := filepath.Join("../db/seeding", f)
		if err := seeding.RunSeeding(ctx, pg, path); err != nil {
			log.Fatalf("seed %s: %v", f, err)
		}
		log.Printf("seeded %s", f)
	}
}
