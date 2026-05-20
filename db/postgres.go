package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type PostgresDB struct {
	Queries *Queries
	Cleanup func()
}

func setupDatabase(ctx context.Context) (*Queries, func()) {
	fmt.Println("establishing database connection...", os.Getenv("LOV_DATABASE_URL"))
	conn, err := pgx.Connect(ctx, os.Getenv("LOV_DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	fmt.Println("database connection established")

	queries := New(conn)
	cleanup := func() {
		conn.Close(ctx)
	}
	return queries, cleanup
}

func NewPostgresDB(ctx context.Context) *PostgresDB {
	queries, cleanup := setupDatabase(ctx)
	return &PostgresDB{
		Queries: queries,
		Cleanup: cleanup,
	}
}
