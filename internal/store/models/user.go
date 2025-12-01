package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Default

	Name         string
	ConfirmedAt  *time.Time
	Email        string  `gorm:"not null;unique"`
	Phone        *string `gorm:"unique"`
	LastSignInAt *time.Time
	Password     string `gorm:"not null"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) IsConfirmed() bool {
	return u.ConfirmedAt != nil
}

type UserProfile struct {
	ID    uuid.UUID
	Name  string
	Email string
	Phone *string
}
