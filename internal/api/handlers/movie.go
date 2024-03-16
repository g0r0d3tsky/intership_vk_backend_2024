package handlers

import (
	"cinema_service/internal/api/handlers/models"
	"cinema_service/internal/domain"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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

// CreateMovieHandler @Summary Create Movie
// @Description Creates a new movie
// @Tags Movies
// @Accept json
// @Param movie body models.Movie true "Movie object"
// @Success 201 "Movie created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Failed to create movie"
// @Router /movies [post]
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

// UpdateMovieHandler @Summary Update Movie
// @Description Updates an existing movie
// @Tags Movies
// @Accept json
// @Param movie_id query string true "Movie ID"
// @Param movie body models.Movie true "Movie object"
// @Success 200 {string}  string "Movie updated successfully"
// @Failure 400 {string} string "Invalid movie ID" or "Invalid request payload"
// @Failure 500 {string} string "Failed to update movie"
// @Router /movies [put]
func (h *MovieHandler) UpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieIDStr := r.URL.Query().Get("id")
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

// DeleteMovieHandler
// @Summary Delete Movie
// @Description Deletes a movie
// @Tags Movies
// @Param movie_id query string true "Movie ID"
// @Success 200 "Movie deleted successfully"
// @Failure 400 {string} string "Invalid movie ID"
// @Failure 500 {string} string "Failed to delete movie"
// @Router /movies [delete]
func (h *MovieHandler) DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieIDStr := r.URL.Query().Get("id")
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

// GetMoviesFilterHandler retrieves movies based on a filter.
// @Summary Get Movies by Filter
// @Description Retrieves movies based on a filter.
// @Tags Movies
// @Param filter query string true "Filter"
// @Success 200 {array} models.Movie
// @Failure 500 {string} 500 "Failed to encode movies"
// @Failure 500 {string} 500 "Failed to get movies"
// @Router /movies [get]
func (h *MovieHandler) GetMoviesFilterHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")

	movies, err := h.service.GetMoviesFilter(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to get movies", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(movies)
	if err != nil {
		http.Error(w, "Failed to encode movies", http.StatusInternalServerError)
		return
	}
}

// GetMoviesBySnippetHandler
// @Summary Get Movies by Snippet
// @Description Retrieves movies based on a snippet
// @Tags Movies
// @Param snippet query string true "Snippet"
// @Success 200 {array} models.Movie
// @Failure 500 {string} 500 "Failed to encode movies" 
// @Failure 500 {string} 500 "Failed to get movies"
// @Router /movies/snippet [get]
func (h *MovieHandler) GetMoviesBySnippetHandler(w http.ResponseWriter, r *http.Request) {
	snippet := r.URL.Query().Get("snippet")

	movies, err := h.service.GetMoviesBySnippet(r.Context(), snippet)
	if err != nil {
		http.Error(w, "Failed to get movies", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(movies)
	if err != nil {
		http.Error(w, "Failed to encode movies", http.StatusInternalServerError)
		return
	}
}

// TODO: authorization
func (h *MovieHandler) RegisterMovie(mux *http.ServeMux,
	authentication Middleware, authorization Middleware) *http.ServeMux {
	mux.HandleFunc("GET /api/v1/movies/filter/{filter}", authentication(h.GetMoviesFilterHandler))
	mux.HandleFunc("GET /api/v1/movies/snippet/{snippet}", authentication(h.GetMoviesBySnippetHandler))
	mux.HandleFunc("POST /api/v1/movies/", authentication(authorization(h.CreateMovieHandler)))
	mux.HandleFunc("PUT /api/v1/movies/{id}", authentication(authorization(h.UpdateMovieHandler)))
	mux.HandleFunc("DELETE /api/v1/movies/{id}", authentication(authorization(h.DeleteMovieHandler)))
	return mux
}
