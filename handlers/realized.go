package handlers

import (
	"log/slog"
	"net/http"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/views/realized"

	"github.com/go-chi/chi/v5"
)

func HandleRealized(w http.ResponseWriter, r *http.Request) error {
	thes_data, err := server.MyS.DB.AllRealizedThesis("id", false, r.URL.Query())
	slog.Info("quere", "q", r.URL.Query())
	if err != nil {
		return err
	}
	return Render(w, r, realized.Index(thes_data))
}

func HandleRealizedFiltered(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()
	for key, val := range q {
		if val[0] == "" {
			q.Del(key)
			slog.Info("Filter", "key", key)
			slog.Info("Filter", "val", val)
		}
	}
	r.URL.RawQuery = q.Encode()
	thes_data, err := server.MyS.DB.AllRealizedThesis("id", false, r.URL.Query())
	if err != nil {
		return err
	}
	return Render(w, r, realized.Results(thes_data))
}
func HandleRealizedDetails(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HRDetails", "id_param", id_param)
	thes_data, err := server.MyS.DB.RealizedThesisByID(id_param)
	slog.Info("quere", "q", r.URL.Query())
	if err != nil {
		return err
	}
	slog.Info("HRealizedDetails", "thes", thes_data)
	return Render(w, r, realized.Details(thes_data))
}

func HandleRealizedEntry(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HREntry", "id_param", id_param)
	thes_data, err := server.MyS.DB.RealizedThesisByID(id_param)
	slog.Info("quere", "q", r.URL.Query())
	if err != nil {
		return err
	}
	slog.Info("HRealizedEntry", "thes", thes_data)
	return Render(w, r, realized.Entry(thes_data))
}

func HandleRealizedGetNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, realized.NewEntry())
}

func HandleRealizedClearNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, realized.EmptySpace())
}
