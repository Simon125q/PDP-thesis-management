package handlers

import (
	"net/http"
	"thesis-management-app/views/components"
)

func HandleLanguageValidation(w http.ResponseWriter, r *http.Request) error {
	lang := r.FormValue("thesisLanguage")
	if lang != "polski" && lang != "angielski" {
		return Render(w, r, components.ErrorMsg("Niepoprawny język", "language-error"))
	}
	return Render(w, r, components.ErrorMsg("", "language-error"))
}
