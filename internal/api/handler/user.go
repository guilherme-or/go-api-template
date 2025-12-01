package handler

import (
	"net/http"

	"github.com/guilherme-or/go-api-template/internal/api"
	"github.com/guilherme-or/go-api-template/internal/domain/user"
)

func ProfileHandler(userService *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, _ := api.GetClaims(r)

		profile, err := userService.GetUserProfile(claims.UserID)
		if err != nil {
			api.RespondJSONErr(w, http.StatusInternalServerError, err)
			return
		}

		api.RespondJSON(w, http.StatusOK, profile)
	}
}
