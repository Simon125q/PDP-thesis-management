package handlers

import (
	"net/http"
	"thesis-management-app/views/realized"
)

func HandleRealized(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, realized.Index())
}
