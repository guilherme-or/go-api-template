package store

import (
	"github.com/guilherme-or/go-api-template/internal/store/models"
	"gorm.io/gorm"
)

type AuthStore interface {
	GetUserByEmail(email string) (*models.User, error)
}

type gormAuthStore struct {
	db *gorm.DB
}

func NewGormAuthStore(db *gorm.DB) AuthStore {
	return &gormAuthStore{db: db}
}

func (s *gormAuthStore) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
