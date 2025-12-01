package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/guilherme-or/go-api-template/config/env"
)

type (
	// Key type for storing claims in context
	ClaimsKey string

	// Custom JWT claims structure
	Claims struct {
		UserID uuid.UUID `json:"user_id"`
		jwt.RegisteredClaims
	}
)

const (
	// Key used to store JWT claims in context
	ClaimsContextKey ClaimsKey = "JWT_CLAIMS"

	// Duration for which the token is valid
	tokenDuration = time.Hour * 1
)

// Creates a JWT token for the given user ID
func GenerateJWTToken(ID uuid.UUID) (string, *Claims, error) {
	now := time.Now()
	claims := Claims{
		UserID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenDuration)),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(env.G.JWTSecret))
	if err != nil {
		return "", nil, err
	}

	return signedToken, &claims, nil
}

// Parses and validates a JWT token string and returns the claims
func ParseJWTToken(t string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(t, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.G.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, jwt.ErrTokenInvalidClaims
	}
}
