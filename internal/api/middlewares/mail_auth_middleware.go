package middlewares

import (
	"anemone_notes/internal/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const AddressContextKey = contextKey("address")

func AuthMailMiddleware(jwtManager *utils.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}

			claims, err := jwtManager.Verify(parts[1])
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AddressContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
