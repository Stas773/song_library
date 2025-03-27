package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func New(logger *logrus.Logger, dsn string) (pool *pgxpool.Pool, err error) {
	const (
		maxAttempts = 10
		delay       = 2 * time.Second
		timeout     = 5 * time.Second
	)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		pool, err = pgxpool.New(ctx, dsn)
		cancel()
		if err == nil {
			logger.Info("successfully connected to PostgreSQL")
			return pool, nil
		}
		logger.Errorf("attempt %d/%d: failed to connect to PostgreSQL: %v", attempt, maxAttempts, err)
		if attempt < maxAttempts {
			logger.Infof("reconnecting in %v...", delay)
			time.Sleep(delay)
		}
	}
	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d attempts: %v", maxAttempts, err)
}
