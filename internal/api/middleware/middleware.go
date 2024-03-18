package middleware

import (
	"cinema_service/internal/api/handlers"
	"cinema_service/internal/domain"
	"cinema_service/internal/usecase"
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

const (
	userCtx = "user_info"
)

//go:generate mockgen -source=middleware.go -destination=mocks/mock.go
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
			handlers.NewErrorResponse(w, http.StatusUnauthorized, "empty auth header")
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			handlers.NewErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		if len(headerParts[1]) == 0 || len(strings.Split(headerParts[1], "")) == 0 {
			handlers.NewErrorResponse(w, http.StatusUnauthorized, "token is empty")
			return
		}

		userInfo, err := m.service.ParseToken(headerParts[1])
		if err != nil {
			handlers.NewErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, userInfo)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
func (m *UserMiddleware) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(userCtx).(*usecase.UserInfo)
		if !ok {
			handlers.NewErrorResponse(w, http.StatusUnauthorized, "User context not found")
			return
		}
		if user.Role != domain.ADMIN {
			handlers.NewErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (m *UserMiddleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mw := NewResponseWriter(w)
		next.ServeHTTP(mw, r)
		statusCode := mw.StatusCode()
		responseSize := mw.Size()

		slog.Info(
			"Method: %s, Path: %s, Status: %s, Size: %d bytes",
			r.Method,
			r.URL.Path,
			strconv.Itoa(statusCode),
			responseSize,
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
