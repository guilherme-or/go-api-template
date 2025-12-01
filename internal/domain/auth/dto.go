package auth

type EmailDTO struct {
	Email string `json:"email"`
}

type LoginMethodDTO struct {
	Email    string  `json:"email"`
	Password *string `json:"password,omitempty"`
	Method   string  `json:"method"`
}

type VerifyMethodDTO struct {
	Valid   bool     `json:"valid"`
	Methods []string `json:"methods"`
}

type AuthenticationDTO struct {
	AccessToken string  `json:"access_token"`
	Claims      *Claims `json:"claims"`
}
