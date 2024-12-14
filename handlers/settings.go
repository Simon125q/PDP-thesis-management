package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/pkgs/validators"
	"thesis-management-app/types"
	"thesis-management-app/views/settings"

	"github.com/go-chi/chi/v5"
)

func HandleSettings(w http.ResponseWriter, r *http.Request) error {
	empl_data, err := server.MyS.DB.AllUniversityEmployee()
	slog.Info("quere", "q", r.URL.Query())
	if err != nil {
		return err
	}
	return Render(w, r, settings.Index(empl_data))
}

func HandleSettingsEntry(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HSEntry", "id_param", id_param)
	empl_data, err := server.MyS.DB.EmployeeEntryByID(id_param)
	slog.Info("quere", "q", r.URL.Query())
	if err != nil {
		return err
	}
	slog.Info("HSettingsEntry", "empl", empl_data)
	return Render(w, r, settings.Entry(empl_data))
}

func HandleSettingsDetails(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HSDetails", "id_param", id_param)
	empl_data, err := server.MyS.DB.EmployeeById(id_param)
	empl_t_count, err := server.MyS.DB.ThesisCountByEmpId(id_param)
	empl_data.ThesisCount = empl_t_count
	slog.Info("quere", "q", r.URL.Query())
	if err != nil {
		return err
	}
	slog.Info("HSettingsDetails", "empl", empl_data)
	return Render(w, r, settings.Details(empl_data, types.UniversityEmployeeErrors{}))
}

func extractEmployeeFromForm(r *http.Request) *types.UniversityEmployee {
	return &types.UniversityEmployee{
		FirstName:            r.FormValue("first_name"),
		LastName:             r.FormValue("last_name"),
		CurrentAcademicTitle: r.FormValue("current_academic_title"),
		DepartmentUnit:       r.FormValue("department_unit"),
	}
}

func HandleSettingsNew(w http.ResponseWriter, r *http.Request) error {
	empl := *extractEmployeeFromForm(r)

	errors, ok := validators.ValidateEmployee(empl)
	if !ok {
		errors.Correct = false
		return Render(w, r, settings.NewEntrySwap(types.UniversityEmployee{}, empl, errors))
	}

	slog.Info("add employee", "correct", true)

	emplId, err := server.MyS.DB.InsertUniversityEmployee(empl)
	if err != nil {
		slog.Error("employee to db", "err", err)
		errors.Correct = false
		return Render(w, r, settings.NewEntrySwap(types.UniversityEmployee{}, empl, errors))
	}

	// Log the successful insertion
	slog.Info("employee to db", "new_id", emplId)

	errors.Correct = true

	// Render the employee page with success
	return Render(w, r, settings.NewEntrySwap(types.UniversityEmployee{}, empl, errors))
}

func HandleSettingsUpdate(w http.ResponseWriter, r *http.Request) error {
	slog.Info("UPDATE", "here", true)

	id_param := chi.URLParam(r, "id")
	slog.Info("UPDATE", "id_param", id_param)

	empl := *extractEmployeeFromForm(r)

	errors, ok := validators.ValidateEmployee(empl)
	if !ok {
		errors.Correct = false
		return Render(w, r, settings.Details(empl, errors))
	}

	var err error
	empl.Id, err = strconv.Atoi(id_param)
	if err != nil {
		slog.Error("Update", "err", err)
		return Render(w, r, settings.Details(empl, errors))
	}

	err = server.MyS.DB.UpdateEmployee(&empl)
	if err != nil {
		slog.Error("Update Employee", "err", err)
		return Render(w, r, settings.Details(empl, errors))
	}

	slog.Info("update employee", "correct", true)
	errors.Correct = true

	return Render(w, r, settings.Details(empl, types.UniversityEmployeeErrors{}))
}

func HandleSettingsGetNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, settings.NewEntry(types.UniversityEmployee{}, types.UniversityEmployeeErrors{}))
}

func HandleSettingsClearNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, settings.EmptySpace())
}
