package auth

import "golang.org/x/crypto/bcrypt"

// Generates a bcrypt hash of the given password
func GeneratePasswordHash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// Compares a bcrypt hashed password with its possible plaintext equivalent
func ComparePasswordHash(hash, plain string) bool {
	if plain == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
