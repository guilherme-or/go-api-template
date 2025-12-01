package handler

import (
	"errors"
	"net/http"

	"github.com/guilherme-or/go-api-template/internal/api"
	"github.com/guilherme-or/go-api-template/internal/domain/auth"
)

func VerifyHandler(authService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := api.ParseJSON[auth.EmailDTO](r)
		if err != nil {
			api.RespondJSONErr(w, http.StatusBadRequest, errors.New("invalid request payload, must have email field"))
			return
		}

		dto, err := authService.VerifyAuthMethods(req.Email)
		if err != nil {
			api.RespondJSONErr(w, http.StatusBadRequest, err)
			return
		}

		api.RespondJSON(w, http.StatusOK, dto)
	}
}

func LoginHandler(authService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := api.ParseJSON[auth.LoginMethodDTO](r)
		if err != nil {
			api.RespondJSONErr(w, http.StatusBadRequest, errors.New("invalid request payload, must have email, password and methods fields"))
			return
		}

		dto, err := authService.AuthenticationGateway(req.Email, req.Method, req.Password)
		if err != nil {
			api.RespondJSONErr(w, http.StatusUnauthorized, err)
			return
		}

		api.RespondJSON(w, http.StatusOK, dto)
	}
}
