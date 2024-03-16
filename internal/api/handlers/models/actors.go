package models

import (
	"cinema_service/internal/domain"
	"time"
)

type Actor struct {
	Name      string    `json:"name"`
	Surname   string    `json:"surname,omitempty"`
	Sex       string    `json:"sex,omitempty"`
	Birthdate time.Time `json:"birthdate,omitempty"`
}

type ActorMovies struct {
	Actor  *domain.Actor
	Movies []*domain.Movie
}
