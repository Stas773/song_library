package storages

import (
	"context"
	"song_library/entities"
)

type Database interface {
	GetSongs(filter entities.Song, page, limit int) ([]entities.Song, error)
	GetSongText(ctx context.Context, id, page, limit int) ([]string, error)
	DeleteSong(ctx context.Context, id int) error
	UpdateSong(ctx context.Context, song *entities.Song) error
	AddSong(ctx context.Context, song *entities.Song) error
}
