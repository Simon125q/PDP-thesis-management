package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/types"
	"thesis-management-app/views/ongoing"

	"github.com/go-chi/chi/v5"
)

func HandleOngoing(w http.ResponseWriter, r *http.Request) error {
	slog.Info("HandleOngoing", "entered", true)
	thes_data, err := server.MyS.DB.AllOngoingThesisEntries("thesis_id", true, 1, PageLimit, r.URL.Query())
	if err != nil {
		slog.Error("HandleOngoing", "err", err)
		return err
	}
	return Render(w, r, ongoing.Index(thes_data))
}

func HandleOngoingDetails(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HandleOngoingDetails", "id_param", id_param)
	thes_data, err := server.MyS.DB.OngoingThesisEntryByID(id_param)
	slog.Info("HandleOngoingDetails", "q", r.URL.Query())
	if err != nil {
		slog.Error("HandleOngoingDetails", "err", err)
		return fmt.Errorf("HandleOngoingDetails -> %v", err)
	}
	return Render(w, r, ongoing.Details(thes_data, types.OngoingThesisEntryErrors{}))
}

func HandleOngoingEntry(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HandleOngoingEntry", "id_param", id_param)
	thes_data, err := server.MyS.DB.OngoingThesisEntryByID(id_param)
	if err != nil {
		return err
	}
	return Render(w, r, ongoing.Entry(thes_data))
}

func HandleOngoingGetNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, ongoing.NewEntry(types.OngoingThesisEntry{}, types.OngoingThesisEntryErrors{}))
}

func HandleOngoingClearNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, ongoing.EmptySpace())
}
