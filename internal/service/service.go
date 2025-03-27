package service

import (
	"context"
	"song_library/entities"
	"song_library/internal/storages"
)

type Service struct {
	db storages.Database
}

func NewService(db storages.Database) *Service {
	return &Service{db: db}
}

func (s *Service) GetSongs(filter entities.Song, page, limit int) ([]entities.Song, error) {
	return s.db.GetSongs(filter, page, limit)
}

func (s *Service) GetSongText(ctx context.Context, id, page, limit int) ([]string, error) {
	return s.db.GetSongText(ctx, id, page, limit)
}

func (s *Service) DeleteSong(ctx context.Context, id int) error {
	return s.db.DeleteSong(ctx, id)
}

func (s *Service) UpdateSong(ctx context.Context, song *entities.Song) error {
	return s.db.UpdateSong(ctx, song)
}

func (s *Service) AddSong(ctx context.Context, song *entities.Song) error {
	return s.db.AddSong(ctx, song)
}
