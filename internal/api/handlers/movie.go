package handlers

import (
	"cinema_service/internal/api/handlers/models"
	"cinema_service/internal/domain"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

//go:generate mockgen -source=movie.go -destination=mocks/movieServiceMock.go

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

// CreateMovieHandler creates a new movie.
// @Summary Create Movie
// @Description Creates a new movie
// @Tags Movies
// @Accept json
// @Security ApiKeyAuth
// @Param movie body models.Movie true "Movie object"
// @Success 201 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /movies [post]
func (h *MovieHandler) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Movie
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
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
		newErrorResponse(w, http.StatusInternalServerError, "Failed to create movie")
		return
	}

	sendJSONResponse(w, http.StatusCreated, statusResponse{
		Status: "Movie created successfully",
	})
}

// UpdateMovieHandler @Summary Update Movie
// @Description Updates an existing movie
// @Tags Movies
// @Accept json
// @Security ApiKeyAuth
// @Param id query string true "Movie ID"
// @Param movie body models.Movie true "Movie object"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /movies [put]
func (h *MovieHandler) UpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieIDStr := r.URL.Query().Get("id")
	if movieIDStr == "" {
		newErrorResponse(w, http.StatusBadRequest, "Movie ID parameter is required")
		return
	}

	id, err := uuid.Parse(movieIDStr)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	var input models.Movie
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	movie := &domain.Movie{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		Date:        input.Date,
		Rating:      input.Rating,
	}

	err = h.service.UpdateMovie(r.Context(), movie)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, errorResponse{
			Message: "Failed to update movie",
		})
		return
	}

	sendJSONResponse(w, http.StatusOK, statusResponse{
		Status: "Movie updated successfully",
	})
}

// DeleteMovieHandler
// @Summary Delete Movie
// @Description Deletes a movie
// @Tags Movies
// @Security ApiKeyAuth
// @Param id query string true "Movie ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /movies [delete]
func (h *MovieHandler) DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieIDStr := r.URL.Query().Get("id")
	if movieIDStr == "" {
		newErrorResponse(w, http.StatusBadRequest, "Movie ID parameter is required")
		return
	}

	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	err = h.service.DeleteMovie(r.Context(), movieID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to delete movie")
		return
	}

	sendJSONResponse(w, http.StatusOK, statusResponse{
		Status: "Movie deleted successfully",
	})
}

// GetMoviesFilterHandler retrieves movies based on a filter.
// @Summary Get Movies by Filter
// @Description Retrieves movies based on a filter.
// @Tags Movies
// @Security ApiKeyAuth
// @Param filter query string true "Filter"
// @Success 200 {array} models.Movie
// @Failure 500 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /movies/filter [get]
func (h *MovieHandler) GetMoviesFilterHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")

	if filter == "" {
		newErrorResponse(w, http.StatusBadRequest, "Filter parameter is required")
		return
	}

	movies, err := h.service.GetMoviesFilter(r.Context(), filter)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to get movies")
		return
	}

	sendJSONResponse(w, http.StatusInternalServerError, movies)

}

// GetMoviesBySnippetHandler retrieves movies based on a snippet.
// @Summary Get Movies by Snippet
// @Description Retrieves movies based on a snippet
// @Tags Movies
// @Security ApiKeyAuth
// @Param snippet query string true "Snippet"
// @Success 200 {array} models.Movie
// @Failure 500 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /movies/snippet [get]
func (h *MovieHandler) GetMoviesBySnippetHandler(w http.ResponseWriter, r *http.Request) {
	snippet := r.URL.Query().Get("snippet")

	if snippet == "" {
		newErrorResponse(w, http.StatusBadRequest, "Snippet parameter is required")
		return
	}

	movies, err := h.service.GetMoviesBySnippet(r.Context(), snippet)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to get movies")
		return
	}

	err = json.NewEncoder(w).Encode(movies)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError,  "Failed to encode movies",
		)
		return
	}
}

// TODO: authorization
func (h *MovieHandler) RegisterMovie(mux *http.ServeMux,
	authentication Middleware, authorization Middleware) *http.ServeMux {
	mux.HandleFunc("GET /api/v1/movies/filter", authentication(h.GetMoviesFilterHandler))
	mux.HandleFunc("GET /api/v1/movies/snippet", authentication(h.GetMoviesBySnippetHandler))
	mux.HandleFunc("POST /api/v1/movies", authentication(authorization(h.CreateMovieHandler)))
	mux.HandleFunc("PUT /api/v1/movies", authentication(authorization(h.UpdateMovieHandler)))
	mux.HandleFunc("DELETE /api/v1/movies", authentication(authorization(h.DeleteMovieHandler)))
	return mux
}
