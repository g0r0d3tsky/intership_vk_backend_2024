package models

import "time"

type Actor struct {
	Name      string    `json:"name"`
	Surname   string    `json:"surname,omitempty"`
	Sex       string    `json:"sex,omitempty"`
	Birthdate time.Time `json:"birthdate,omitempty"`
}
