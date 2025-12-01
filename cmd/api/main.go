package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/guilherme-or/go-api-template/config"
	"github.com/guilherme-or/go-api-template/config/env"
	"github.com/guilherme-or/go-api-template/internal/api"
	"github.com/guilherme-or/go-api-template/internal/api/handler"
	"github.com/guilherme-or/go-api-template/internal/domain/auth"
	"github.com/guilherme-or/go-api-template/internal/domain/user"
	"github.com/guilherme-or/go-api-template/internal/store"
	"github.com/guilherme-or/go-api-template/internal/store/database"
)

func main() {
	config.MustSetup()

	slog.Debug("Loaded environment", "environment", env.G)

	start(wire())
}

func wire() *api.Server {
	// connections
	db, err := database.NewGORMPostgresConnection(env.G.Database.DSN, env.G.Debug)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
	}

	// store implementations
	authStore := store.NewGormAuthStore(db)
	userStore := store.NewGormUserStore(db)

	// services
	authService := auth.NewAuthService(authStore)
	userService := user.NewUserService(userStore)

	// router
	r := api.NewRouter()
	r.Register(func(m *chi.Mux) {
		m.Use(api.JSONMiddleware)

		auth := chi.NewRouter()
		auth.Post("/verify", handler.VerifyHandler(authService))
		auth.Post("/login", handler.LoginHandler(authService))

		profile := chi.NewRouter()
		profile.Use(api.JWTMiddleware)
		profile.Get("/", handler.ProfileHandler(userService))

		m.Mount("/auth", auth)
		m.Mount("/profile", profile)
	})

	// server
	return api.NewServer(env.G.Port, r.Handler())
}

func start(s *api.Server) {
	s.Start()
	defer s.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
