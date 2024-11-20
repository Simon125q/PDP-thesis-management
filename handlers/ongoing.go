package handlers

import (
	"net/http"
	"thesis-management-app/views/ongoing"
)

func HandleOngoing(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, ongoing.Index())
}

func HandleOngoingGetNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, ongoing.NewEntry())
}
