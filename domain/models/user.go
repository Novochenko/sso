package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	HashPassword []byte
}

func (u User) ValidateRegister() error {
	return validation.ValidateStruct(
		u,
		validation.Field(u.Email, validation.Required, is.Email),
		validation.Field(u.HashPassword, validation.NilOrNotEmpty),
	)
}
