package models

import (
	"time"
)

type Movie struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date,omitempty" format:"2006-01-02"`
	Rating      float32   `json:"rating,omitempty"`
}
