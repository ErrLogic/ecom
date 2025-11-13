package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ErrLogic/ecom/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config *config
}

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hi"))

		if err != nil {
			return
		}
	})

	productService := products.NewService()
	productHandler := products.NewHandlers(productService)
	r.Get("/products", productHandler.ListProducts)

	return r
}

// run
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Listening on %s\n", srv.Addr)

	return srv.ListenAndServe()
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

func loadConfig() *config {
	return &config{
		addr: os.Getenv("SERVER_ADDR"),
		db: dbConfig{
			dsn: os.Getenv("DATABASE_DSN"),
		},
	}
}
