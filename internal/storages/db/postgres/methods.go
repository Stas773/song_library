package postgres

import (
	"context"
	"fmt"
	"song_library/entities"
	"strings"
)

type PostgresStorage struct {
	db Client
}

func NewPostgresStorage(db Client) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) GetSongs(filter entities.Song, page, limit int) ([]entities.Song, error) {
	var songs []entities.Song
	str := "SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE group_name ILIKE $1 AND song_name ILIKE $2 AND release_date ILIKE $3 AND link ILIKE $4 LIMIT $5 OFFSET $6"
	rows, err := s.db.Query(context.Background(), str,
		"%"+filter.GroupName+"%",
		"%"+filter.SongName+"%",
		"%"+filter.ReleaseDate+"%",
		"%"+filter.Link+"%",
		limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var song entities.Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return songs, nil
}

func (s *PostgresStorage) GetSongText(ctx context.Context, id int, page, limit int) ([]string, error) {
	var text string
	err := s.db.QueryRow(ctx, "SELECT text FROM songs WHERE id = $1", id).Scan(&text)
	if err != nil {
		return nil, err
	}

	verses := strings.Split(text, "\n\n")
	fmt.Println(verses)
	start := (page - 1) * limit
	if start >= len(verses) {
		return []string{}, nil
	}

	end := start + limit
	if end > len(verses) {
		end = len(verses)
	}
	fmt.Println(verses[start:end])
	fmt.Println(strings.Join(verses[start:end], "\n\n"))

	return verses[start:end], nil
}

func (s *PostgresStorage) DeleteSong(ctx context.Context, id int) error {
	_, err := s.db.Exec(ctx, "DELETE FROM songs WHERE id = $1", id)
	return err
}

func (s *PostgresStorage) UpdateSong(ctx context.Context, song *entities.Song) error {
	_, err := s.db.Exec(ctx, "UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5 WHERE id = $6",
		song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link, song.ID)
	return err
}

func (s *PostgresStorage) AddSong(ctx context.Context, song *entities.Song) error {
	_, err := s.db.Exec(ctx, "INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)",
		song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link)
	return err
}
