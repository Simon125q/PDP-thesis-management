package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"thesis-management-app/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	router := chi.NewMux()
	router.Use(middleware.Logger)
	router.Use(handlers.WithUser)

	router.Handle("/*", public())
	router.Get("/", handlers.Make(handlers.HandleHome))
	router.Get("/login", handlers.Make(handlers.HandleLogin))
	router.Get("/ongoing", handlers.Make(handlers.HandleOngoing))
	router.Get("/realized", handlers.Make(handlers.HandleRealized))

	listenAddr := os.Getenv("LISTEN_ADDR")

	slog.Info("HTTP server", "ListenAddr", listenAddr)
	http.ListenAndServe(listenAddr, router)
}
