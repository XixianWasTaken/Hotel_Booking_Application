package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"learningGo/cmd/internal/config"
	"testing"
)

func TestRouts(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//pass
	default:
		t.Error(fmt.Sprintf("type is not *chi mux, type is %T", v))
	}
}
