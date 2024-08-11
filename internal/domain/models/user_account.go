package models

import "github.com/google/uuid"

type UserAccount struct {
	UserId             uuid.UUID
	UserName           string
	ProfilePicturePath string
}
