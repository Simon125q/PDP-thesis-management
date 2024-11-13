package handlers

import (
	"net/http"
	"thesis-management-app/views/settings"
)

func HandleSettings(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, settings.Index())
}
