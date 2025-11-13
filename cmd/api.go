package main

import (
	"log"
	"net/http"
	"os"
	"time"

	repo "github.com/ErrLogic/ecom/internal/adapters/postgresql/sqlc"
	"github.com/ErrLogic/ecom/internal/orders"
	"github.com/ErrLogic/ecom/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config *config
	db     *pgx.Conn
}

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hi"))

		if err != nil {
			return
		}
	})

	// Products
	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandlers(productService)
	r.Get("/products", productHandler.ListProducts)

	// Order
	orderService := orders.NewService(repo.New(app.db), app.db)
	orderHandler := orders.NewHandler(orderService)
	r.Post("/orders", orderHandler.PlaceOrder)

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
			dsn: os.Getenv("GOOSE_DBSTRING"),
		},
	}
}
