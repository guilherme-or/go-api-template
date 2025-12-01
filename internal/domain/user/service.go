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

type UserService struct {
	store store.UserStore
}

func NewUserService(store store.UserStore) *UserService {
	return &UserService{store: store}
}

func (s *UserService) GetUserProfile(ID uuid.UUID) (*ProfileDTO, error) {
	profile, err := s.store.GetUserProfile(ID)
	if err != nil {
		slog.Error("Unable to retrieve user profile for UUID "+ID.String(), "error", err, "uuid", ID)
		return nil, ErrUserProfile
	}

	dto := ProfileDTO{}
	dto.PopulateModel(profile)
	return &dto, nil
}
