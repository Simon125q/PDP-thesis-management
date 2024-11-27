package handlers

import (
	"log/slog"
	"net/http"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/views/settings"
)

func HandleSettings(w http.ResponseWriter, r *http.Request) error {
	empl_data, err := server.MyS.DB.AllUniversityEmployee()
	slog.Info("quere", "q", r.URL.Query())
	if err != nil {
		return err
	}
	return Render(w, r, settings.Index(empl_data))
}
