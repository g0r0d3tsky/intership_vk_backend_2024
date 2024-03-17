package middleware

import (
	"cinema_service/internal/domain"
	"cinema_service/internal/usecase"
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
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

func (m *UserMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
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

		if len(headerParts[1]) == 0 || len(strings.Split(headerParts[1], "")) == 0 {
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
func (m *UserMiddleware) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
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
func (m *UserMiddleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		mw := NewResponseWriter(w)
		next.ServeHTTP(mw, r)
		statusCode := mw.StatusCode()
		responseSize := mw.Size()

		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)

		slog.Info(
			"Method: %s, Path: %s, Status: %s, Size: %d bytes, Time: %s",
			r.Method,
			r.URL.Path,
			strconv.Itoa(statusCode),
			responseSize,
			elapsedTime,
		)
	})
}

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK, 0}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.size += size
	return size, err
}

func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

func (rw *ResponseWriter) Size() int {
	return rw.size
}
