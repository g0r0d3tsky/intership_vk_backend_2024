package domain

import (
	"github.com/google/uuid"
	"time"
)

type Actor struct {
	ID        uuid.UUID
	Name      string
	Surname   string
	Sex       string
	Birthdate time.Time
}
