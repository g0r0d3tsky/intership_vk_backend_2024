package domain

import (
	"github.com/google/uuid"
	"time"
)

type Movie struct {
	ID          uuid.UUID
	Title       string
	Description string
	Date        time.Time
	Rating      int
}
