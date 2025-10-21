package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mightyfzeus/stage-one/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr   string
	apiUrl string
	db     dbConfig
}

type dbConfig struct {
	dbAddr       string
	maxOpenConns int
	maxIdleTime  string
	maxIdleConns int
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Public middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/strings", app.CreateStringController)
	r.Get("/strings/{string_value}", app.GetStringHandler)
	r.Delete("/strings/{string_value}", app.DeleteStringHandler)
	r.Get("/strings", app.GetAllStringsHandler)
	r.Get("/strings/filter-by-natural-language", app.FilterByNaturalLanguageHandler)

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	return srv.ListenAndServe()
}
