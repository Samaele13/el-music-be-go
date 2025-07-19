package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Song struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	ImageURL string `json:"imageUrl"`
	SongURL  string `json:"songUrl"`
}

type PostgresStore struct {
	Db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=elmusic password=supersecret dbname=elmusic_dev sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return &PostgresStore{Db: db}, nil
}

func (s *PostgresStore) GetRecentlyPlayed() ([]Song, error) {
	rows, err := s.Db.Query("SELECT id, title, artist, image_url, song_url FROM songs WHERE section = 'recently_played'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.ImageURL, &song.SongURL); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}
