package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Playlist struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

// ... (Struct Song, Category, User biarkan sama)
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
	IsVerified   bool
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

func (s *PostgresStore) GetUserPlaylists(userID string) ([]Playlist, error) {
	rows, err := s.Db.Query("SELECT id, name, owner_id FROM playlists WHERE owner_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []Playlist
	for rows.Next() {
		var p Playlist
		if err := rows.Scan(&p.ID, &p.Name, &p.OwnerID); err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}
	return playlists, nil
}

func (s *PostgresStore) CreatePlaylist(name, ownerID string) (*Playlist, error) {
	var newPlaylist Playlist
	err := s.Db.QueryRow(
		"INSERT INTO playlists (name, owner_id) VALUES ($1, $2) RETURNING id, name, owner_id",
		name, ownerID,
	).Scan(&newPlaylist.ID, &newPlaylist.Name, &newPlaylist.OwnerID)

	if err != nil {
		return nil, err
	}
	return &newPlaylist, nil
}

// ... (sisa fungsi User, Song, Category biarkan sama)
func (s *PostgresStore) CreateUser(name, email, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	verificationToken := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	_, err = s.Db.Exec(
		"INSERT INTO users (name, email, password_hash, verification_token, verification_token_expires_at) VALUES ($1, $2, $3, $4, $5)",
		name, email, string(hashedPassword), verificationToken, expiresAt,
	)
	if err != nil {
		return "", err
	}
	return verificationToken, nil
}

func (s *PostgresStore) VerifyUser(token string) error {
	res, err := s.Db.Exec(
		"UPDATE users SET is_verified = true, verification_token = NULL, verification_token_expires_at = NULL WHERE verification_token = $1 AND verification_token_expires_at > NOW()",
		token,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *PostgresStore) GetUserByEmail(email string) (*User, error) {
	var user User
	err := s.Db.QueryRow("SELECT id, name, email, password_hash, is_verified FROM users WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.IsVerified)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgresStore) SetPasswordResetToken(email string) (string, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour) // Token berlaku 1 jam
	_, err := s.Db.Exec(
		"UPDATE users SET reset_password_token = $1, reset_password_token_expires_at = $2 WHERE email = $3",
		token, expiresAt, email,
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *PostgresStore) ResetPassword(token, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	res, err := s.Db.Exec(
		"UPDATE users SET password_hash = $1, reset_password_token = NULL, reset_password_token_expires_at = NULL WHERE reset_password_token = $2 AND reset_password_token_expires_at > NOW()",
		string(hashedPassword), token,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Token tidak valid atau kedaluwarsa
	}
	return nil
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
