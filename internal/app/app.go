package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"song_library/internal/handlers"
	"song_library/internal/storages/db/postgres"
	"time"

	_ "song_library/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(logger *logrus.Logger, dsn string, port string, songHandlers *handlers.SongHandler) {

	pool, err := postgres.New(logger, dsn)
	if err != nil {
		logger.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pool.Close()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/api/v1/songs", songHandlers.GetSongs)
	router.GET("/api/v1/songs/:id/text", songHandlers.GetSongText)
	router.POST("/api/v1/songs", songHandlers.AddSong)
	router.PUT("/api/v1/songs/:id", songHandlers.UpdateSong)
	router.DELETE("/api/v1/songs/:id", songHandlers.DeleteSong)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server didn't started %s\n", err)
			os.Exit(1)
		}
	}()
	logger.Infof("Server started on port:%s\n", port)

	<-done
	logger.Info("Stop signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Shutdown error %s\n", err)
	}
	logger.Info("Server stoped")
}
