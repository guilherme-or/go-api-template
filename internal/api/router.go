package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	mux *chi.Mux
}

// Creates configured server router with external mux
func NewRouter() *Router {
	return &Router{mux: chi.NewRouter()}
}

func (ro *Router) Handler() http.Handler {
	return ro.mux
}

func (ro *Router) Register(routes func(*chi.Mux)) {
	// global middlewares
	ro.mux.Use(middleware.Logger)
	ro.mux.Use(middleware.Recoverer)
	ro.mux.Use(middleware.RealIP)
	ro.mux.Use(middleware.StripSlashes)
	ro.mux.Use(middleware.Heartbeat("/ping"))

	// custom handlers
	ro.mux.NotFound(ro.notFoundHandler)
	ro.mux.MethodNotAllowed(ro.methodNotAllowedHandler)

	routes(ro.mux)
}

func (ro *Router) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	RespondJSONErr(w, http.StatusNotFound, errors.New("resource not found"))
}

func (ro *Router) methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	RespondJSONErr(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
}
