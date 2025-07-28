package handler

import (
	"el-music-be/internal/database"
	"encoding/json"
	"net/http"
)

type SearchHandler struct {
	Store *database.PostgresStore
}

func NewSearchHandler(store *database.PostgresStore) *SearchHandler {
	return &SearchHandler{Store: store}
}

func (h *SearchHandler) HandleSearchSongs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		// Return empty list if query is empty
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]database.Song{})
		return
	}

	songs, err := h.Store.SearchSongs(query)
	if err != nil {
		http.Error(w, "Failed to search songs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}
