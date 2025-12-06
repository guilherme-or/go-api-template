package store

import (
	"github.com/google/uuid"
	"github.com/guilherme-or/go-api-template/internal/store/models"
	"gorm.io/gorm"
)

type UserStore interface {
	GetUserByID(id uuid.UUID) (*models.User, error)
}

type gormUserStore struct {
	db *gorm.DB
}

func NewGormUserStore(db *gorm.DB) UserStore {
	return &gormUserStore{db: db}
}

func (s *gormUserStore) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User

	result := s.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
