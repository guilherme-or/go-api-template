package user

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilherme-or/go-api-template/internal/store"
)

var (
	ErrUserProfile = errors.New("unable to retrieve user profile")
)

type UserService interface {
	GetUserProfile(ID uuid.UUID) (*ProfileDTO, error)
}

type userServiceImpl struct {
	store store.UserStore
}

func NewUserService(store store.UserStore) UserService {
	return &userServiceImpl{store: store}
}

func (s *userServiceImpl) GetUserProfile(ID uuid.UUID) (*ProfileDTO, error) {
	profile, err := s.store.GetUserByID(ID)
	if err != nil {
		slog.Error("Unable to retrieve user", "error", err, "uuid", ID)
		return nil, ErrUserProfile
	}

	dto := ProfileDTO{}
	dto.PopulateModel(profile)
	return &dto, nil
}
