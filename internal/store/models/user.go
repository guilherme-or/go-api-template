package models

import (
	"time"
)

var RoleUser = "user"
var RoleAdmin = "admin"

type User struct {
	Default

	Name         string
	Email        string  `gorm:"not null;unique"`
	Phone        *string `gorm:"unique"`
	Role         string  `gorm:"not null"`
	LastSignInAt *time.Time
	Password     string `gorm:"not null"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) Roles() []string {
	roles := []string{}
	switch u.Role {
	case RoleAdmin:
		roles = append(roles, RoleAdmin)
		fallthrough
	case RoleUser:
		roles = append(roles, RoleUser)
	}
	return roles
}
