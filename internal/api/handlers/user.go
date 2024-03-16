package handlers

import (
	"cinema_service/internal/usecase"
	"context"
	"encoding/json"
	"net/http"
)

type UserService interface {
	GenerateToken(ctx context.Context, login string, password string) (string, error)
	ParseToken(token string) (*usecase.UserInfo, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

type signInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Sign In
// @Description Authenticates a user and returns a token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body signInInput true "Sign In Input"
// @Success 200 {object} map[string]interface{} "Token"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /signin [post]
func (s *UserHandler) signIn(w http.ResponseWriter, r *http.Request) {
	var input signInInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := s.service.GenerateToken(r.Context(), input.Login, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (h *UserHandler) RegisterUser(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("POST /api/v1/users/", h.signIn)
	return mux
}
