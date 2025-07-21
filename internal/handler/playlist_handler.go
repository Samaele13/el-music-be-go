package handler

import (
	"el-music-be/internal/database"
	"el-music-be/internal/middleware"
	"encoding/json"
	"net/http"

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
