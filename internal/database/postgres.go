package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Song struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	ImageURL string `json:"imageUrl"`
	SongURL  string `json:"songUrl"`
}

type Category struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
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

func (s *PostgresStore) CreateUser(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.Db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", name, email, string(hashedPassword))
	return err
}

func (s *PostgresStore) GetUserByEmail(email string) (*User, error) {
	var user User
	err := s.Db.QueryRow("SELECT id, name, email, password_hash FROM users WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
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

func (s *PostgresStore) GetMadeForYou() ([]Song, error) {
	rows, err := s.Db.Query("SELECT id, title, artist, image_url, song_url FROM songs WHERE section = 'made_for_you'")
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

func (s *PostgresStore) GetSearchCategories() ([]Category, error) {
	rows, err := s.Db.Query("SELECT id, name, image_url FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.ImageURL); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
