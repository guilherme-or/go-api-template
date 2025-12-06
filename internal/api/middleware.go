package api

import (
	"context"
	"errors"
	"mime"
	"net/http"
	"strings"

	"github.com/guilherme-or/go-api-template/internal/domain/auth"
)

var (
	ContentTypeHeader = "Content-Type"
	JSONContentType   = mime.TypeByExtension(".json")

	ErrNotJSONContentType      = errors.New("unsupported Content-Type header, expected application/json")
	ErrMissingBearerToken      = errors.New("missing bearer token in Authorization header")
	ErrInvalidBearerToken      = errors.New("invalid or expired bearer token")
	ErrInsufficientPermissions = errors.New("user has insufficient permissions to access this resource")
)

// Ensures that client and server communicate using JSON via Content-Type header
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(ContentTypeHeader, JSONContentType)

		if r.ContentLength == 0 {
			next.ServeHTTP(w, r)
			return
		}

		clientContentType := strings.TrimSpace(r.Header.Get(ContentTypeHeader))

		if clientContentType != JSONContentType {
			RespondJSONErr(w, http.StatusUnsupportedMediaType, ErrNotJSONContentType)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Validates JWT token from Authorization header and populates request context with claims
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerContent := strings.TrimSpace(r.Header.Get("Authorization"))
		if headerContent == "" {
			RespondJSONErr(w, http.StatusUnauthorized, ErrMissingBearerToken)
			return
		}

		parts := strings.Split(headerContent, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			RespondJSONErr(w, http.StatusUnauthorized, ErrMissingBearerToken)
			return
		}

		claims, err := auth.ParseJWTToken(parts[1])
		if err != nil {
			RespondJSONErr(w, http.StatusUnauthorized, ErrInvalidBearerToken)
			return
		}

		// Insert claims in context for further handlers to use
		ctx := r.Context()
		ctx = context.WithValue(ctx, auth.ClaimsContextKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.Claims)
			if !ok || claims == nil {
				RespondJSONErr(w, http.StatusUnauthorized, ErrInvalidBearerToken)
				return
			}

			for _, userRole := range claims.Roles {
				for _, allowedRole := range allowedRoles {
					if userRole == allowedRole {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			RespondJSONErr(w, http.StatusForbidden, ErrInsufficientPermissions)
		})
	}
}
