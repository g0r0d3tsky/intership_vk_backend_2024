package domain

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	ADMIN = "ADMIN"
	USER  = "USER"
)

type User struct {
	ID        uuid.UUID
	Role      string
	Login     string
	Password  []byte
	CreatedAt time.Time
}

func GeneratePasswordHash(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
func (u *User) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	u.Password = hash
	return nil
}

func (u *User) Matches(plaintextPassword string) (bool, error) {
	p, err := GeneratePasswordHash(plaintextPassword)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(p, u.Password)
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
