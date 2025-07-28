package middleware

import (
	"context"
	"el-music-be/internal/auth"
	"el-music-be/internal/database"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"
const IsSubscribedKey contextKey = "isSubscribed"

func JWTMiddleware(store *database.PostgresStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			claims := &auth.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return auth.JwtKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			user, err := store.GetUserByEmail(claims.Subject)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			isSubscribed := user.SubscriptionStatus == "active" && (user.SubscriptionExpiresAt.Valid && user.SubscriptionExpiresAt.Time.After(time.Now()))

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, IsSubscribedKey, isSubscribed)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
