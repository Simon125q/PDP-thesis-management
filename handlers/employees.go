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

func HandleEmployees(w http.ResponseWriter, r *http.Request) error {
	empl_data, err := server.MyS.DB.AllUniversityEmployeeEntries("employee_id", true, r.URL.Query())
	slog.Info("quere", "q", r.URL.Query())
	if err != nil {
		return err
	}
	return Render(w, r, settings.ResultsUsers(empl_data))
}

func HandleEmployeesEntry(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HSEntry", "id_param", id_param)
	empl_data, err := server.MyS.DB.EmployeeEntryByID(id_param)
	slog.Info("quere", "q", r.URL.Query())
	if err != nil {
		return err
	}
	slog.Info("HSettingsEntry", "empl", empl_data)
	return Render(w, r, settings.Entry_Empl(empl_data))
}

func HandleEmployeesDetails(w http.ResponseWriter, r *http.Request) error {
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
	return Render(w, r, settings.Details_Empl(empl_data, types.UniversityEmployeeEntryErrors{}))
}

func extractEmployeeFromForm(r *http.Request) *types.UniversityEmployeeEntry {
	return &types.UniversityEmployeeEntry{
		FirstName:            r.FormValue("first_name"),
		LastName:             r.FormValue("last_name"),
		CurrentAcademicTitle: r.FormValue("current_academic_title"),
		DepartmentUnit:       r.FormValue("department_unit"),
	}
}

func HandleEmployeesNew(w http.ResponseWriter, r *http.Request) error {
	empl := *extractEmployeeFromForm(r)

	errors, ok := validators.ValidateEmployee(empl)
	if !ok {
		errors.Correct = false
		return Render(w, r, settings.NewEntrySwap_Empl(types.UniversityEmployeeEntry{}, empl, errors))
	}

	slog.Info("add employee", "correct", true)

	emplId, err := server.MyS.DB.InsertUniversityEmployee(empl)
	if err != nil {
		slog.Error("employee to db", "err", err)
		errors.Correct = false
		return Render(w, r, settings.NewEntrySwap_Empl(types.UniversityEmployeeEntry{}, empl, errors))
	}

	slog.Info("employee to db", "new_id", emplId)
	empl.Id = int(emplId)
	errors.Correct = true

	return Render(w, r, settings.NewEntrySwap_Empl(empl, types.UniversityEmployeeEntry{}, errors))
}

func HandleEmployeesUpdate(w http.ResponseWriter, r *http.Request) error {
	slog.Info("UPDATE", "here", true)

	id_param := chi.URLParam(r, "id")
	slog.Info("UPDATE", "id_param", id_param)

	empl := *extractEmployeeFromForm(r)

	errors, ok := validators.ValidateEmployee(empl)
	if !ok {
		errors.Correct = false
		return Render(w, r, settings.Details_Empl(empl, errors))
	}

	var err error
	empl.Id, err = strconv.Atoi(id_param)
	if err != nil {
		slog.Error("Update", "err", err)
		return Render(w, r, settings.Details_Empl(empl, errors))
	}

	err = server.MyS.DB.UpdateEmployee(&empl)
	if err != nil {
		slog.Error("Update Employee", "err", err)
		return Render(w, r, settings.Details_Empl(empl, errors))
	}

	slog.Info("update employee", "correct", true)
	errors.Correct = true

	return Render(w, r, settings.Details_Empl(empl, types.UniversityEmployeeEntryErrors{}))
}

func HandleEmployeesGetNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, settings.NewEntry_Empl(types.UniversityEmployeeEntry{}, types.UniversityEmployeeEntryErrors{}))
}

func HandleEmployeesClearNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, settings.EmptySpace_Empl())
}
