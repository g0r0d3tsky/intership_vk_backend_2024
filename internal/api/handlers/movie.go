package handlers

import (
	"cinema_service/internal/api/handlers/models"
	"cinema_service/internal/domain"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

type MovieService interface {
	CreateMovie(ctx context.Context, movie *domain.Movie) error
	UpdateMovie(ctx context.Context, movie *domain.Movie) error
	DeleteMovie(ctx context.Context, movieID uuid.UUID) error
	GetMoviesFilter(ctx context.Context, filter string) ([]*domain.Movie, error)
	GetMoviesBySnippet(ctx context.Context, snippet string) ([]*domain.Movie, error)
}

type MovieHandler struct {
	service MovieService
}

func NewMovieHandler(service MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

func (h *MovieHandler) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Movie
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	movie := &domain.Movie{
		Title:       input.Title,
		Description: input.Description,
		Date:        input.Date,
		Rating:      input.Rating,
	}
	err = h.service.CreateMovie(r.Context(), movie)
	if err != nil {
		http.Error(w, "Failed to create movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *MovieHandler) UpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieIDStr := r.URL.Query().Get("movie_id")
	_, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	var input models.Movie
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	movie := &domain.Movie{
		Title:       input.Title,
		Description: input.Description,
		Date:        input.Date,
		Rating:      input.Rating,
	}
	err = h.service.UpdateMovie(r.Context(), movie)
	if err != nil {
		http.Error(w, "Failed to update movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MovieHandler) DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieIDStr := r.URL.Query().Get("movie_id")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteMovie(r.Context(), movieID)
	if err != nil {
		http.Error(w, "Failed to delete movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MovieHandler) GetMoviesFilterHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")

	movies, err := h.service.GetMoviesFilter(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to get movies", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(movies)
	if err != nil {
		return
	}
}

func (h *MovieHandler) GetMoviesBySnippetHandler(w http.ResponseWriter, r *http.Request) {
	snippet := r.URL.Query().Get("snippet")

	movies, err := h.service.GetMoviesBySnippet(r.Context(), snippet)
	if err != nil {
		http.Error(w, "Failed to get movies", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(movies)
	if err != nil {
		return
	}
}
