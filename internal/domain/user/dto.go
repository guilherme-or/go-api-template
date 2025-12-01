package user

import (
	"github.com/google/uuid"
	"github.com/guilherme-or/go-api-template/internal/store/models"
)

type ProfileDTO struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Phone *string   `json:"phone,omitempty"`
}

func (dto *ProfileDTO) PopulateModel(model *models.UserProfile) {
	dto.ID = model.ID
	dto.Name = model.Name
	dto.Email = model.Email
	dto.Phone = model.Phone
}
