package main

import (
	"net/http"

	"github.com/ChunHou23/booking-service/pkg/config"
	"github.com/ChunHou23/booking-service/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/AboutMe", handlers.Repo.AboutMe)

	return mux
}
