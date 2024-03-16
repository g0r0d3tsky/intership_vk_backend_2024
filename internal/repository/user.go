package repository

import (
	"cinema_service/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StorageUser struct {
	db *pgxpool.Pool
}

func NewUserStorage(dbPool *pgxpool.Pool) StorageUser {
	StorageUser := StorageUser{
		db: dbPool,
	}
	return StorageUser
}
func (s *StorageUser) GetUser(ctx context.Context, login string, password string) (*domain.User, error) {
	user := &domain.User{}
	pass, err := domain.GeneratePasswordHash(password)
	if err != nil {
		return nil, err
	}
	if err = s.db.QueryRow(
		ctx,
		`SELECT id, login, password, role, created_at FROM "users" u WHERE u.login = $1 AND u.password = $2`, login, pass,
	).Scan(&user.ID, &user.Login, &user.Password, &user.Role, &user.CreatedAt); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

func (s *StorageUser) CreateUser(ctx context.Context, u *domain.User) error {
	u.ID = uuid.New()
	u.CreatedAt = time.Now()
	if _, err := s.db.Exec(ctx,
		`INSERT INTO "users" (id, login, password, role, created_at) 
			VALUES ($1, $2, $3, $4, $5)`,
		&u.ID, &u.Login, &u.Password, &u.Role, &u.CreatedAt,
	); err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_login_key"`:
			return ErrDuplicateLogin
		default:
			return err
		}
	}
	return nil
}
func (s *StorageUser) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if _, err := s.db.Exec(ctx,
		`DELETE FROM "users" WHERE id=$1`,
		userID,
	); err != nil {
		return err
	}
	return nil
}
