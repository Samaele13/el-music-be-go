package handler

import (
	"el-music-be/internal/database"
	"encoding/json"
	"net/http"
)

type SongHandler struct {
	Store *database.PostgresStore
}

func NewSongHandler(store *database.PostgresStore) *SongHandler {
	return &SongHandler{Store: store}
}

func (h *SongHandler) HandleGetRecentlyPlayed(w http.ResponseWriter, r *http.Request) {
	songs, err := h.Store.GetRecentlyPlayed()
	if err != nil {
		http.Error(w, "Failed to fetch songs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (h *SongHandler) HandleGetMadeForYou(w http.ResponseWriter, r *http.Request) {
	songs, err := h.Store.GetMadeForYou()
	if err != nil {
		http.Error(w, "Failed to fetch songs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (h *SongHandler) HandleGetSearchCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Store.GetSearchCategories()
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
