package seeding

import (
	"context"
	"fmt"
	"os"

	"lov/db"
)

func RunSeeding(ctx context.Context, pg *db.PostgresDB, filePath string) error {
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", filePath, err)
	}
	if err := pg.RawExec(ctx, string(sqlBytes)); err != nil {
		return fmt.Errorf("exec %s: %w", filePath, err)
	}
	return nil
}
