package auth

import (
	"errors"
	"log/slog"
	"net/mail"

	"github.com/guilherme-or/go-api-template/internal/store"
)

var (
	AuthMethodPassword = "password"

	ErrInvalidEmail          = errors.New("invalid email address")
	ErrMethodVerification    = errors.New("could not verify authentication method for provided email address")
	ErrUnsupportedMethod     = errors.New("given authentication method is not supported")
	ErrRequiredPassword      = errors.New("password is required for password authentication method")
	ErrInvalidUserOrPassword = errors.New("invalid user or password")
	ErrTokenGeneration       = errors.New("error generating authentication token")
)

type AuthService struct {
	store store.AuthStore
}

func NewAuthService(store store.AuthStore) *AuthService {
	return &AuthService{store: store}
}

func (s *AuthService) validateEmailAddress(email string) error {
	if email == "" {
		slog.Error("Empty email address", "email", email)
		return ErrInvalidEmail
	}

	if _, err := mail.ParseAddress(email); err != nil {
		slog.Error("Error parsing provided email address", "email", email, "error", err)
		return ErrInvalidEmail
	} else {
		return nil
	}
}

func (s *AuthService) VerifyAuthMethods(email string) (*VerifyMethodDTO, error) {
	dto := &VerifyMethodDTO{
		Valid:   false,
		Methods: []string{},
	}

	if err := s.validateEmailAddress(email); err != nil {
		return dto, err
	}

	_, err := s.store.GetUserByEmail(email)
	if err != nil {
		slog.Error("Error finding user with provided email address", "email", email, "error", err)
		return dto, ErrMethodVerification
	}

	dto.Valid = true

	// If the user exists, we assume they have password authentication enabled
	dto.Methods = append(dto.Methods, AuthMethodPassword)

	return dto, nil
}

func (s *AuthService) AuthenticationGateway(email, method string, password *string) (*AuthenticationDTO, error) {
	switch method {
	case AuthMethodPassword:
		if password == nil {
			slog.Error(ErrRequiredPassword.Error())
			return nil, ErrRequiredPassword
		}
		return s.AuthenticateWithPassword(email, *password)
	default:
		slog.Error(ErrUnsupportedMethod.Error(), "method", method)
		return nil, ErrUnsupportedMethod
	}
}

func (s *AuthService) AuthenticateWithPassword(email, password string) (*AuthenticationDTO, error) {
	if err := s.validateEmailAddress(email); err != nil {
		return nil, err
	}

	user, err := s.store.GetUserByEmail(email)
	if err != nil {
		slog.Error("Error finding user with provided email address", "email", email, "error", err)
		return nil, ErrInvalidUserOrPassword
	}

	if !ComparePasswordHash(user.Password, password) {
		slog.Error("Password check failed", "email", email)
		return nil, ErrInvalidUserOrPassword
	}

	token, claims, err := GenerateJWTToken(user)
	if err != nil {
		slog.Error(ErrTokenGeneration.Error(), "email", email, "error", err, "token", token)
		return nil, ErrTokenGeneration
	}

	return &AuthenticationDTO{AccessToken: token, Claims: claims}, nil
}
