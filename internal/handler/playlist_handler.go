package handler

import (
	"el-music-be/internal/database"
	"el-music-be/internal/middleware"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type PlaylistHandler struct {
	Store *database.PostgresStore
}

func NewPlaylistHandler(store *database.PostgresStore) *PlaylistHandler {
	return &PlaylistHandler{Store: store}
}

type CreatePlaylistRequest struct {
	Name string `json:"name"`
}

type AddSongRequest struct {
	SongID string `json:"song_id"`
}

func (h *PlaylistHandler) HandleRemoveSongFromPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not get user ID from context", http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	playlistID := vars["playlistId"]
	songID := vars["songId"]

	err := h.Store.RemoveSongFromPlaylist(playlistID, songID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "user does not own this playlist") {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to remove song", http.StatusNotFound)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Song removed from playlist successfully"})
}

func (h *PlaylistHandler) HandleAddSongToPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not get user ID from context", http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	playlistID := vars["id"]
	var req AddSongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := h.Store.AddSongToPlaylist(playlistID, req.SongID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "user does not own this playlist") {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Song already in playlist", http.StatusConflict)
		} else {
			http.Error(w, "Failed to add song to playlist", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Song added to playlist successfully"})
}

func (h *PlaylistHandler) HandleGetUserPlaylists(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not get user ID from context", http.StatusInternalServerError)
		return
	}
	playlists, err := h.Store.GetUserPlaylists(userID)
	if err != nil {
		http.Error(w, "Failed to fetch playlists", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlists)
}

func (h *PlaylistHandler) HandleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not get user ID from context", http.StatusInternalServerError)
		return
	}
	var req CreatePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	playlist, err := h.Store.CreatePlaylist(req.Name, userID)
	if err != nil {
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(playlist)
}

func (h *PlaylistHandler) HandleGetPlaylistByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not get user ID from context", http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	playlistID := vars["id"]
	playlist, err := h.Store.GetPlaylistByID(playlistID, userID)
	if err != nil {
		http.Error(w, "Playlist not found or access denied", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlist)
}
