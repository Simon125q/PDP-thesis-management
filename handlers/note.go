package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/views/components"

	"github.com/go-chi/chi/v5"
)

func HandleNote(w http.ResponseWriter, r *http.Request) error {
	realized_id, err := strconv.Atoi(chi.URLParam(r, "realized_id"))
	if err != nil {
		slog.Info("HandleNote", "err", err)
	}
	ongoing_id, err := strconv.Atoi(chi.URLParam(r, "ongoing_id"))
	if err != nil {
		slog.Info("HandleNote", "err", err)
	}
	user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		slog.Info("HandleNote", "err", err)
	}
	note, err := server.MyS.DB.GetNote(realized_id, ongoing_id, user_id)
	if err != nil {
		return fmt.Errorf("HandleNote error getting note %v", err)
	}
	return Render(w, r, components.Note(note))
}
