package repository

import (
	"cinema_service/internal/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type StorageActor struct {
	db *pgxpool.Pool
}

func NewStorageActor(dbPool *pgxpool.Pool) StorageActor {
	StorageActor := StorageActor{
		db: dbPool,
	}
	return StorageActor
}

func (s *StorageActor) CreateActor(ctx context.Context, act *domain.Actor) error {
	act.ID = uuid.New()
	if _, err := s.db.Exec(ctx,
		`INSERT INTO "actors" (id, name, surname, sex, birthdate) 
			VALUES ($1, $2, $3, $4, $5)`,
		&act.ID, &act.Name, &act.Surname, &act.Sex, &act.Birthdate,
	); err != nil {
		return fmt.Errorf("create actor: %w", err)
	}
	return nil
}
func (s *StorageActor) UpdateActor(ctx context.Context, act *domain.Actor) error {
	if _, err := s.db.Exec(
		ctx,
		`UPDATE "actors" SET name = $2, surname = $3, sex = $4, birthdate = $5
              WHERE id = $1`,
		&act.ID, &act.Name, &act.Surname, &act.Sex, &act.Birthdate,
	); err != nil {
		return fmt.Errorf("update actor: %w", err)
	}
	return nil
}
func (s *StorageActor) GetActors(ctx context.Context) (map[*domain.Actor][]*domain.Movie, error) {
	var actors []*domain.Actor
	rows, err := s.db.Query(ctx, `
		SELECT a.id, a.name, a.surname, a.sex, a.birthdate, m.id, m.title, m.description, m.rating, m.created_at 
		FROM actors a
		INNER JOIN actors_movies am ON a.id = am.actors_movie_id
		INNER JOIN movies m ON am.movies_id = m.id
	`)
	if err != nil {
		return nil, fmt.Errorf("get actors: %w", err)
	}
	defer rows.Close()
	actorFilms := make(map[*domain.Actor][]*domain.Movie)
	for rows.Next() {
		actor := &domain.Actor{}
		movie := &domain.Movie{}

		if err = rows.Scan(
			&actor.ID, &actor.Name, &actor.Surname, &actor.Sex, &actor.Birthdate,
			&movie.ID, &movie.Title, &movie.Description, &movie.Rating, &movie.Date,
		); err != nil {
			return nil, fmt.Errorf("get actors: %w", err)
		}

		var existingActor *domain.Actor
		for _, a := range actors {
			if a.ID == actor.ID {
				existingActor = a
				break
			}
		}

		if existingActor == nil {
			existingActor = actor
			actors = append(actors, existingActor)
		}
		actorFilms[existingActor] = append(actorFilms[existingActor], movie)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("get actors: %w", err)
	}

	return actorFilms, nil
}
func (s *StorageActor) DeleteActor(ctx context.Context, actorID uuid.UUID) error {
	result, err := s.db.Exec(ctx,
		`DELETE FROM "movies" WHERE id=$1`,
		actorID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to delete actor")
	}

	rowsAffected := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("actor not found")
	}

	return nil
}