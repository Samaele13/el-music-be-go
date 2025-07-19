package main

import (
	"el-music-be/internal/database"
	"el-music-be/internal/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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

func main() {
	store, err := database.NewPostgresStore()
	if err != nil {
		log.Fatal("Could not connect to the database: ", err)
	}

	songHandler := handler.NewSongHandler(store)

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/songs/recently-played", songHandler.HandleGetRecentlyPlayed).Methods("GET")

	handler := corsMiddleware(r)

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
