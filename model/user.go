package model

import "time"

type User struct {
	Id           int    `json:"id"`
	Email        string `json:"email" validate:"required,email"`
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
	CreationDate time.Time
}
