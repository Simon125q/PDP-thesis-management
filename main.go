package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"thesis-management-app/handlers"
	"thesis-management-app/pkgs/ldap"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/types/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	server.MyS.Router = chi.NewRouter()
	ldap.SetupLDAP()

	db, err := sql.Open("sqlite3", "./diploma_database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server.MyS.DB = &sqlite.Model{DB: db}

	server.MyS.Router.Use(middleware.Logger)
	server.MyS.Router.Use(handlers.WithUser)
	server.MyS.Router.Use(handlers.RefreshSession)

	server.MyS.Router.Group(func(r chi.Router) {
		r.Handle("/*", public())
		r.Get("/", handlers.Make(handlers.HandleHome))
		r.Get("/login", handlers.Make(handlers.HandleLogin))
		r.Post("/login", handlers.Make(handlers.HandleLoginPost))
		r.Post("/logout", handlers.Make(handlers.HandleLogoutPost))
	})

	server.MyS.Router.Group(func(r chi.Router) {
		r.Use(handlers.WithAuth)
		r.Get("/ongoing", handlers.Make(handlers.HandleOngoing))
		r.Get("/ongoing/new", handlers.Make(handlers.HandleOngoingGetNew))
		r.Get("/ongoing/clear-new", handlers.Make(handlers.HandleOngoingClearNew))
		r.Get("/realized", handlers.Make(handlers.HandleRealized))
		r.Get("/realized/generate_excel", handlers.Make(handlers.HandleRealizedGenerateExcel))
		r.Get("/realized/filter", handlers.Make(handlers.HandleRealizedFiltered))
		r.Get("/realized/clear-new", handlers.Make(handlers.HandleRealizedClearNew))
		r.Get("/realized/{id}", handlers.Make(handlers.HandleRealizedEntry))
		r.Get("/realized/details/{id}", handlers.Make(handlers.HandleRealizedDetails))
	})

	server.MyS.Router.Group(func(r chi.Router) {
		r.Use(handlers.WithAdminRights)
		r.Get("/settings", handlers.Make(handlers.HandleSettings))
		r.Get("/realized/new", handlers.Make(handlers.HandleRealizedGetNew))
		r.Post("/realized", handlers.Make(handlers.HandleRealizedNew))
	})

	listenAddr := os.Getenv("LISTEN_ADDR")

	slog.Info("HTTP server", "ListenAddr", listenAddr)
	http.ListenAndServe(listenAddr, server.MyS.Router)
}
