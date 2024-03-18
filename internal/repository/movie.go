package repository

import (
	"cinema_service/internal/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StorageMovie struct {
	db *pgxpool.Pool
}

func NewStorageMovie(dbPool *pgxpool.Pool) StorageMovie {
	StorageMovie := StorageMovie{
		db: dbPool,
	}
	return StorageMovie
}

func (s *StorageMovie) GetMovieByID(ctx context.Context, movieID uuid.UUID) (*domain.Movie, error) {
	movie := &domain.Movie{}

	if err := s.db.QueryRow(
		ctx,
		`SELECT id, title, description, rating, created_at FROM "movies" u WHERE u.id = $1`, movieID,
	).Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Rating, &movie.Rating, &movie.Date); err != nil {
		return nil, fmt.Errorf("get movie by id: %w", err)
	}

	return movie, nil
}

func (s *StorageMovie) CreateMovie(ctx context.Context, movie *domain.Movie) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO "movies" (id, title, description, rating, created_at) VALUES($1, $2, $3, $4, $5)`,
		&movie.ID, &movie.Title, &movie.Description, &movie.Rating, &movie.Rating, &movie.Date,
	)
	if err != nil {
		return fmt.Errorf("create movie: %w", err)
	}
	return nil
}
func (s *StorageMovie) GetMovies(ctx context.Context) ([]*domain.Movie, error) {
	var movies []*domain.Movie
	rows, err := s.db.Query(
		ctx,
		`SELECT id, title, description, rating, created_at FROM movies`)
	if err != nil {
		return nil, fmt.Errorf("get movies: %w", err)
	}
	defer rows.Close() 

	for rows.Next() {
		movie := &domain.Movie{}

		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Rating, &movie.Date); err != nil {
			return nil, fmt.Errorf("get movies: %w", err)
		}

		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil { 
		return nil, fmt.Errorf("get movies: %w", err)
	}

	return movies, nil
}
func (s *StorageMovie) GetMoviesBySnippet(ctx context.Context, snippet string) ([]*domain.Movie, error) {
	var movies []*domain.Movie
	rows, err := s.db.Query(
		ctx,
		`SELECT movies.id, movies.title, movies.description, movies.rating, movies.created_at
		FROM movies
		JOIN actors_movies ON movies.id = actors_movies.movies_id
		JOIN actors ON actors_movies.actors_movie_id = actors.id
		WHERE movies.title LIKE '%' || $1 || '%'
		OR actors.name LIKE '%' || $1 || '%'`,
		snippet)
	if err != nil {
		return nil, fmt.Errorf("get movie by snippet: %w", err)
	}
	for rows.Next() {
		movie := &domain.Movie{}
		if err = rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Rating, &movie.Date); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (s *StorageMovie) UpdateMovie(ctx context.Context, movie *domain.Movie) error {
	if _, err := s.db.Exec(
		ctx,
		`UPDATE "movies" SET title = $2, description = $3, rating = $4, created_at = $5
		WHERE id = $1`,
		&movie.ID, &movie.Title, &movie.Description, &movie.Rating, &movie.Date,
	); err != nil {
		return fmt.Errorf("update movie: %w", err)
	}
	return nil
}
func (s *StorageMovie) DeleteMovie(ctx context.Context, movieID uuid.UUID) error {
	if _, err := s.db.Exec(ctx,
		`DELETE FROM "movies" WHERE id=$1`,
		movieID,
	); err != nil {
		return fmt.Errorf("delete movie: %w", err)
	}
	return nil
}
