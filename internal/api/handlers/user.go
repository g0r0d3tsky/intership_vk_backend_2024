package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

//go:generate mockgen -source=user.go -destination=mocks/userServiceMock.go

type UserService interface {
	GenerateToken(ctx context.Context, login string, password string) (string, error)
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

type signInResponse struct {
	Token string `json:"token"`
}

// @Summary Sign In
// @Description Authenticates a user and returns a token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param sigIn body signInInput true "Sign In Input"
// @Success 200 {object} signInResponse "Token response"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /signIn [post]
func (s *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input signInInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, "Unmarshalling error")
		return
	}

	token, err := s.service.GenerateToken(r.Context(), input.Login, input.Password)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, "Generating Token error")
		return
	}

	response := signInResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) RegisterUser(mux *http.ServeMux, logging Middleware) *http.ServeMux {
	mux.HandleFunc("POST /api/v1/signIn", logging(h.SignIn))
	return mux
}
