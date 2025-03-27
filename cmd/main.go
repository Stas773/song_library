package main

import (
	"fmt"
	"song_library/internal/app"
	"song_library/internal/config"
	"song_library/internal/handlers"
	"song_library/internal/service"
	"song_library/internal/storages/db/migrations"
	"song_library/internal/storages/db/postgres"
	"song_library/pkg/logger"
)

// @title Music Library API
// @version 1.0
// @description This is a sample server for a music library.
// @host localhost:8080
// @BasePath /

func main() {
	config := config.LoadConfig()
	logger := logger.NewLogger()
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	pgxPool, err := postgres.New(logger, dsn)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer pgxPool.Close()

	migrations.Migrate(logger, pgxPool)

	songRepo := postgres.NewPostgresStorage(pgxPool)
	service := service.NewService(songRepo)
	songHandlers := handlers.NewSongHandler(logger, service)

	app.Run(logger, dsn, config.Port, songHandlers)
}
