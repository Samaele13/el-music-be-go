package main

import (
	"el-music-be/internal/database"
	"el-music-be/internal/handler"
	"el-music-be/internal/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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
	authHandler := handler.NewAuthHandler(store)

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	authRoutes := api.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/register", authHandler.HandleRegister).Methods("POST")
	authRoutes.HandleFunc("/login", authHandler.HandleLogin).Methods("POST")
	authRoutes.HandleFunc("/verify", authHandler.HandleVerifyEmail).Methods("GET")

	protectedRoutes := api.PathPrefix("").Subrouter()
	protectedRoutes.Use(middleware.JWTMiddleware)
	protectedRoutes.HandleFunc("/songs/recently-played", songHandler.HandleGetRecentlyPlayed).Methods("GET")
	protectedRoutes.HandleFunc("/songs/made-for-you", songHandler.HandleGetMadeForYou).Methods("GET")
	protectedRoutes.HandleFunc("/categories/search", songHandler.HandleGetSearchCategories).Methods("GET")

	handler := corsMiddleware(r)

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
