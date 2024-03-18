package repository

import (
	"cinema_service/internal/domain"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	if err := s.db.QueryRow(
		ctx,
		`SELECT id, login, password, role, created_at FROM "users" u WHERE u.login = $1`, login,
	).Scan(&user.ID, &user.Login, &user.Password, &user.Role, &user.CreatedAt); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("get user: %w", err)
	}
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return nil, fmt.Errorf("wrong password: %w", err)
	}

	return user, nil
}
