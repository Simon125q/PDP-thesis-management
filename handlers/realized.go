package handlers

import (
	"log/slog"
	"net/http"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/views/realized"
)

func HandleRealized(w http.ResponseWriter, r *http.Request) error {
	thes_data, err := server.MyS.DB.AllRealizedThesis("id", false, r.URL.Query())
	slog.Info("quere", "q", r.URL.Query())
	if err != nil {
		return err
	}
	return Render(w, r, realized.Index(thes_data))
}
