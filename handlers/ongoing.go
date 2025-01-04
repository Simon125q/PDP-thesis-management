package handlers

import (
	"log/slog"
	"net/http"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/types"
	"thesis-management-app/views/ongoing"
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

func HandleOngoingGetNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, ongoing.NewEntry(types.OngoingThesisEntry{}, types.OngoingThesisEntryErrors{}))
}

func HandleOngoingClearNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, ongoing.EmptySpace())
}
