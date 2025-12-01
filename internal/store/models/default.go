package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Default struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}