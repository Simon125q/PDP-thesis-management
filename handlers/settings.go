package handlers

import (
	"log/slog"
	"net/http"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/views/settings"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, settings.Index())
}

func HandleSettingsEmployees(w http.ResponseWriter, r *http.Request) error {
	empl_data, err := server.MyS.DB.AllUniversityEmployeeEntries("employee_id", true, r.URL.Query())
	if err != nil {
		slog.Error("Fetch employees", "err", err)
		return err
	}
	return Render(w, r, settings.ResultsUsers(empl_data))
}

func HandleSettingsCourses(w http.ResponseWriter, r *http.Request) error {
	courses_data, err := server.MyS.DB.AllCourses()
	if err != nil {
		slog.Error("Fetch courses", "err", err)
		return err
	}
	return Render(w, r, settings.ResultsCourses(courses_data))
}

func HandleSettingsSpecializations(w http.ResponseWriter, r *http.Request) error {
	specs_data, err := server.MyS.DB.AllSpecializations()
	if err != nil {
		slog.Error("Fetch specializations", "err", err)
		return err
	}
	return Render(w, r, settings.ResultsSpecs(specs_data))
}
