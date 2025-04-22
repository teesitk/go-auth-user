package http

import (
	"context"
	"encoding/json"
	"go-auth-user/domain"
	"net/http"
	"strings"
)

type key string

const userContextKey key = "user"

func JWTMiddleware(authService domain.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			erJson, _ := json.MarshalIndent(errAuthResp{Code: http.StatusUnauthorized, Error: "Missing or invalid Authorization header"}, "", "    ")
			http.Error(w, string(erJson), http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		user, err := authService.ParseToken(token)
		if err != nil {
			erJson, _ := json.MarshalIndent(errAuthResp{Code: http.StatusUnauthorized, Error: "Invalid token"}, "", "    ")
			http.Error(w, string(erJson), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) *domain.User {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	if !ok {
		return nil
	}
	return user
}
