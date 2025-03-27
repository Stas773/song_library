package migrations

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func Migrate(logger *logrus.Logger, db *pgxpool.Pool) {
	query, err := os.ReadFile("internal/storages/db/migrations/000001_create_songs_table.up.sql")
	if err != nil {
		logger.Errorf("failed to read migration file: %v", err)
	}

	_, err = db.Exec(context.Background(), string(query))
	if err != nil {
		logger.Errorf("failed to execute migration: %v", err)
	}
}
