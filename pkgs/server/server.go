package server

import (
	"thesis-management-app/types/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/store"
)

const ServerContextKey = "server"

type Server struct {
	Router        *chi.Mux
	DB            *sqlite.Model
	Authenticator auth.Authenticator
	Cache         store.Cache
}

var MyS = &Server{}
