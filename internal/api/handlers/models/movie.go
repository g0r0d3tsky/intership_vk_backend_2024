package models

import (
	"time"
)

type Movie struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	Rating      int       `json:"rating,omitempty"`
}
