package main

import (
	"encoding/json"
	"log"
	"net/http"
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func recentlyPlayedHandler(w http.ResponseWriter, r *http.Request) {
	songs := []Song{
		{ID: "recent_1", Title: "Lagu Dari Server 1", Artist: "Artis Go", ImageURL: "https://placehold.co/300x300/5C7E6D/FFFFFF?text=Go+1", SongURL: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3"},
		{ID: "recent_2", Title: "Lagu Dari Server 2", Artist: "Artis Go", ImageURL: "https://placehold.co/300x300/5C7E6D/FFFFFF?text=Go+2", SongURL: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-2.mp3"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func madeForYouHandler(w http.ResponseWriter, r *http.Request) {
	songs := []Song{
		{ID: "mfy_1", Title: "Mix Go Server 1", Artist: "Playlist Go", ImageURL: "https://placehold.co/300x300/1C1C1E/FFFFFF?text=Go+Mix+1", SongURL: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-3.mp3"},
		{ID: "mfy_2", Title: "Mix Go Server 2", Artist: "Playlist Go", ImageURL: "https://placehold.co/300x300/1C1C1E/FFFFFF?text=Go+Mix+2", SongURL: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-4.mp3"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func searchCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories := []Category{
		{ID: "cat_go_1", Name: "Go Pop", ImageURL: "https://placehold.co/300x300/E57373/FFFFFF?text=Go+Pop"},
		{ID: "cat_go_2", Name: "Go Rock", ImageURL: "https://placehold.co/300x300/81C784/FFFFFF?text=Go+Rock"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/songs/recently-played", recentlyPlayedHandler)
	mux.HandleFunc("/api/v1/songs/made-for-you", madeForYouHandler)
	mux.HandleFunc("/api/v1/categories/search", searchCategoriesHandler)

	handler := corsMiddleware(mux)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}