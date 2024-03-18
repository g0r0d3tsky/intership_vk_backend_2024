package usecase

import (
	"cinema_service/internal/domain"
	"context"
	"fmt"
	"sort"

	"github.com/google/uuid"
)

//go:generate mockgen -source=movie.go -destination=mocks/movieMock.go

type MovieRepo interface {
	CreateMovie(ctx context.Context, movie *domain.Movie) error
	GetMovies(ctx context.Context) ([]*domain.Movie, error)
	GetMoviesBySnippet(ctx context.Context, snippet string) ([]*domain.Movie, error)
	UpdateMovie(ctx context.Context, movie *domain.Movie) error
	DeleteMovie(ctx context.Context, movieID uuid.UUID) error
}

type MovieService struct {
	repo MovieRepo
}

func NewMovieService(repo MovieRepo) *MovieService {
	return &MovieService{repo: repo}
}


func (s *MovieService) CreateMovie(ctx context.Context, movie *domain.Movie) error {
	err := s.repo.CreateMovie(ctx, movie)
	if err != nil {
		return fmt.Errorf("create movie: %w", err)
	}
	return nil
}

func (s *MovieService) UpdateMovie(ctx context.Context, movie *domain.Movie) error {
	err := s.repo.UpdateMovie(ctx, movie)
	if err != nil {
		return fmt.Errorf("update movie: %w", err)
	}
	return nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, movieID uuid.UUID) error {
	err := s.repo.DeleteMovie(ctx, movieID)
	if err != nil {
		return fmt.Errorf("delete movie: %w", err)
	}
	return nil
}

func (s *MovieService) GetMovies(ctx context.Context) ([]*domain.Movie, error) {
	movies, err := s.repo.GetMovies(ctx)
	if err != nil {
		return nil, fmt.Errorf("get movies: %w", err)
	}
	return movies, nil
}

func (s *MovieService) GetMoviesFilter(ctx context.Context, filter string) ([]*domain.Movie, error) {
	movies, err := s.repo.GetMovies(ctx)
	if err != nil {
		return nil, fmt.Errorf("get movies: %w", err)
	}

	switch filter {
	case "title":
		sort.Slice(movies, func(i, j int) bool {
			return movies[i].Title < movies[j].Title
		})
	case "rating":
		sort.Slice(movies, func(i, j int) bool {
			return movies[i].Rating > movies[j].Rating
		})
	case "created_at":
		sort.Slice(movies, func(i, j int) bool {
			return movies[i].Date.After(movies[j].Date)
		})
	case "":
		sort.Slice(movies, func(i, j int) bool {
			return movies[i].Rating > movies[j].Rating
		})
	default:
		return nil, fmt.Errorf("invalid filter")
	}

	return movies, nil
}

func (s *MovieService) GetMoviesBySnippet(ctx context.Context, snippet string) ([]*domain.Movie, error) {
	movies, err := s.repo.GetMoviesBySnippet(ctx, snippet)
	if err != nil {
		return nil, fmt.Errorf("get movies by snippet: %w", err)
	}
	return movies, nil
}
