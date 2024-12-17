package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"thesis-management-app/handlers"
	"thesis-management-app/pkgs/ldap"
	"thesis-management-app/pkgs/logging"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/types/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func main() {

	logging.SetupLogger()

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	server.MyS.Router = chi.NewRouter()
	ldap.SetupLDAP()

	db_path := os.Getenv("DB_PATH")
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server.MyS.DB = &sqlite.Model{DB: db}

	// server.MyS.DB.DB.SetConnMaxLifetime(time.Minute)
	// server.MyS.DB.DB.SetMaxOpenConns(1)
	// _, _ = server.MyS.DB.DB.Exec("PRAGMA busy_timeout = 5000") // Log retries
	// _, err = server.MyS.DB.DB.Exec("PRAGMA journal_mode=WAL;")
	// if err != nil {
	// 	log.Fatal("Failed to enable WAL mode:", err)
	// }

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
		r.Get("/realized/generate_excel", handlers.Make(handlers.HandleRealizedGenerateExcel))
		r.Get("/realized/filter", handlers.Make(handlers.HandleRealizedFiltered))
		r.Get("/realized/clear-new", handlers.Make(handlers.HandleRealizedClearNew))
		r.Get("/realized/details/{id}", handlers.Make(handlers.HandleRealizedDetails))
		r.Get("/realized/entry/{id}", handlers.Make(handlers.HandleRealizedEntry))
		r.Get("/realized/autocompleteThesisTitlePolish", handlers.Make(handlers.HandleAutocompleteThesisTitlePolish))
		r.Get("/realized/autocompleteStudentSurname", handlers.Make(handlers.HandleAutocompleteStudentSurname))
		r.Get("/realized/autocompleteStudentNumber", handlers.Make(handlers.HandleAutocompleteStudentNumber))
		r.Get("/realized/autocompleteStudentNameAndSurname", handlers.Make(handlers.HandleAutocompleteStudentNameAndSurname))
		r.Get("/realized/autocompleteSupervisorName", handlers.Make(handlers.HandleAutocompleteSupervisorName))
		r.Get("/realized/autocompleteSupervisorSurname", handlers.Make(handlers.HandleAutocompleteSupervisorSurname))
		r.Get("/realized/autocompleteSupervisorTitle", handlers.Make(handlers.HandleAutocompleteSupervisorTitle))
		r.Get("/realized/autocompleteAssistantSupervisorName", handlers.Make(handlers.HandleAutocompleteAssistantSupervisorName))
		r.Get("/realized/autocompleteAssistantSupervisorSurname", handlers.Make(handlers.HandleAutocompleteAssistantSupervisorSurname))
		r.Get("/realized/autocompleteSupervisorNameAndSurname", handlers.Make(handlers.HandleAutocompleteSupervisorNameAndSurname))
		r.Get("/realized/autocompleteAssistantSupervisorTitle", handlers.Make(handlers.HandleAutocompleteAssistantSupervisorTitle))
		r.Get("/realized/autocompleteReviewerName", handlers.Make(handlers.HandleAutocompleteReviewerName))
		r.Get("/realized/autocompleteReviewerSurname", handlers.Make(handlers.HandleAutocompleteReviewerSurname))
		r.Get("/realized/autocompleteReviewerTitle", handlers.Make(handlers.HandleAutocompleteReviewerTitle))
		r.Get("/realized/autocompleteChairName", handlers.Make(handlers.HandleAutocompleteChairName))
		r.Get("/realized/autocompleteChairSurname", handlers.Make(handlers.HandleAutocompleteChairSurname))
		r.Get("/realized/autocompleteChairTitle", handlers.Make(handlers.HandleAutocompleteChairTitle))
		r.Get("/realized/autocompleteCourse", handlers.Make(handlers.HandleAutocompleteCourse))
		r.Get("/realized", handlers.Make(handlers.HandleRealized))
	})

	server.MyS.Router.Group(func(r chi.Router) {
		r.Use(handlers.WithAdminRights)
		r.Get("/settings", handlers.Make(handlers.HandleSettings))
		r.Get("/settings/clear-new", handlers.Make(handlers.HandleSettingsClearNew))
		r.Get("/settings/{id}", handlers.Make(handlers.HandleSettingsEntry))
		r.Post("/settings", handlers.Make(handlers.HandleSettingsNew))
		r.Get("/settings/new", handlers.Make(handlers.HandleSettingsGetNew))
		r.Get("/settings/details/{id}", handlers.Make(handlers.HandleSettingsDetails))
		r.Put("/settings/{id}", handlers.Make(handlers.HandleSettingsUpdate))
		r.Get("/realized/new", handlers.Make(handlers.HandleRealizedGetNew))
		r.Post("/realized", handlers.Make(handlers.HandleRealizedNew))
		r.Put("/realized/{id}", handlers.Make(handlers.HandleRealizedUpdate))
	})

	listenAddr := os.Getenv("LISTEN_ADDR")

	slog.Info("HTTP server", "ListenAddr", listenAddr)
	http.ListenAndServe(listenAddr, server.MyS.Router)
}
