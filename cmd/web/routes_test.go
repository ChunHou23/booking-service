package main

import (
	"fmt"
	"testing"

	"github.com/ChunHou23/booking-service/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoute(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Error(fmt.Sprintf("Type is not chi.Mux, but is %T", v))
	}
}
