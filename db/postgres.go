package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type PostgresDB struct {
	Queries *Queries
	conn *pgx.Conn
	Cleanup func()
}

func setupDatabase(ctx context.Context) (*Queries, func(), *pgx.Conn) {
	fmt.Println("establishing database connection...")
	conn, err := pgx.Connect(ctx, os.Getenv("LOV_DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	fmt.Println("database connection established")

	queries := New(conn)
	cleanup := func() {
		conn.Close(ctx)
	}
	return queries, cleanup, conn
}

func NewPostgresDB(ctx context.Context) *PostgresDB {
	queries, cleanup, conn := setupDatabase(ctx)
	return &PostgresDB{
		Queries: queries,
		Cleanup: cleanup,
		conn: conn,
	}
}

func (p *PostgresDB) Raw(ctx context.Context, query string, args ...any) (error) {
	_, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (p *PostgresDB) RawExec(ctx context.Context, query string, args ...any) (error) {
	_, err := p.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}
