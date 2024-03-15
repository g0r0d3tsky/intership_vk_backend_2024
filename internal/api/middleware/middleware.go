package middleware

import (
	"cinema_service/internal/domain"
	"cinema_service/internal/usecase"
	"context"
	"net/http"
	"strings"
)

const (
	userCtx = "user_info"
)

type UserService interface {
	GenerateToken(ctx context.Context, login string, password string) (string, error)
	ParseToken(token string) (*usecase.UserInfo, error)
}

type UserMiddleware struct {
	service UserService
}

func NewUserMiddleware(service UserService) *UserMiddleware {
	return &UserMiddleware{service: service}
}

func (m *UserMiddleware) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			http.Error(w, "empty auth header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		if len(headerParts[1]) == 0 {
			http.Error(w, "token is empty", http.StatusUnauthorized)
			return
		}

		userInfo, err := m.service.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, userInfo)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
func (m *UserMiddleware) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(userCtx).(*domain.User)
		if !ok {
			http.Error(w, "User context not found", http.StatusUnauthorized)
			return
		}
		if user.Role != domain.ADMIN {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
