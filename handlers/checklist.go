package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/types"
	"thesis-management-app/views/ongoing"

	"github.com/go-chi/chi/v5"
)

func HandleChecklist(w http.ResponseWriter, r *http.Request) error {
	thesis_id, err := strconv.Atoi(chi.URLParam(r, "thesis_id"))
	if err != nil {
		slog.Error("HandleChecklist", "err", err)
		return err
	}
	slog.Info("HandleChecklist", "thesis_id", thesis_id)
	tasks, err := server.MyS.DB.GetTasksByThesisId(thesis_id)
	if err != nil {
		slog.Error("HandleChecklist", "err", err)
		return err
	}
	return Render(w, r, ongoing.Checklist(tasks, thesis_id))
}

func HandleInsertTask(w http.ResponseWriter, r *http.Request) error {
	ongoingThesisID, _ := strconv.Atoi(chi.URLParam(r, "thesis_id"))
	task := types.Task{
		Content:         r.FormValue("new_task"),
		OngoingThesisID: ongoingThesisID,
	}
	id, err := server.MyS.DB.InsertTask(task)
	if err != nil {
		slog.Error("HandleInsertTask", "err", err)
		return err
	}
	task.Id = int(id)
	return Render(w, r, ongoing.NewTaskSwap(task))
}
