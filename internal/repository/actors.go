package repository

import (
	"cinema_service/internal/domain"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
		//TODO:
		return err
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
		return err
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
		return nil, err
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
			return nil, err
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
		//TODO: how to store films for actors?
		actorFilms[existingActor] = append(actorFilms[existingActor], movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return actorFilms, nil
}
func (s *StorageActor) DeleteActor(ctx context.Context, actorID uuid.UUID) error {
	if _, err := s.db.Exec(ctx,
		`DELETE FROM "movies" WHERE id=$1`,
		actorID,
	); err != nil {
		return err
	}
	return nil
}