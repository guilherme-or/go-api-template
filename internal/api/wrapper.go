package api

import (
	"encoding/json"
	"net/http"

	"github.com/guilherme-or/go-api-template/internal/domain/auth"
)

type ErrResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"error"`
}

// Writes the given body as JSON to the response writer with the specified status code
func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// Writes a JSON error response with the specified status code and error message
func RespondJSONErr(w http.ResponseWriter, code int, err error) {
	RespondJSON(w, code, ErrResponse{
		StatusCode:   code,
		ErrorMessage: err.Error(),
	})
}

// Decodes the JSON body of the request into the provided generic type T
func ParseJSON[T interface{}](r *http.Request) (*T, error) {
	var dst T
	err := json.NewDecoder(r.Body).Decode(&dst)
	return &dst, err
}

func GetClaims(r *http.Request) (*auth.Claims, bool) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.Claims)
	return claims, ok
}
