package handler

import (
	"el-music-be/internal/database"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type LyricsHandler struct {
	Store *database.PostgresStore
}

func NewLyricsHandler(store *database.PostgresStore) *LyricsHandler {
	return &LyricsHandler{Store: store}
}

func (h *LyricsHandler) HandleGetLyrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songID := vars["songId"]

	lyrics, err := h.Store.GetLyricsForSong(songID)
	if err != nil {
		http.Error(w, "Failed to fetch lyrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lyrics)
}
