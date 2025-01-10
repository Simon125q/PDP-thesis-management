package handlers

import (
	"log/slog"
	"net/http"
	"thesis-management-app/pkgs/validators"
	"thesis-management-app/views/components"
)

func HandleStudentNumberValidate(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handle", "url", r.URL.Query())
	q := r.URL.Query()
	index := q.Get("studentNumber")
	err, _ := validators.ValidateIndex(index)
	return Render(w, r, components.ErrorMsg(err))
}

func HandleStudentNameValidate(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handle", "url", r.URL.Query())
	q := r.URL.Query()
	name := q.Get("firstNameStudent")
	err, _ := validators.ValidateName(name)
	return Render(w, r, components.ErrorMsg(err))
}

func HandleStudentSurnameValidate(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handle", "url", r.URL.Query())
	q := r.URL.Query()
	name := q.Get("lastNameStudent")
	err, _ := validators.ValidateName(name)
	return Render(w, r, components.ErrorMsg(err))
}
