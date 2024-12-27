package handlers

import (
	"github.com/go-chi/chi/v5"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/pkgs/validators"
	"thesis-management-app/types"
	"thesis-management-app/views/settings"
)

type SpecFilter struct {
	Search string
	Order  string
	SortBy string
}

func HandleSpecializations(w http.ResponseWriter, r *http.Request) error {
	specializations, err := server.MyS.DB.AllSpecializations()
	if err != nil {
		slog.Error("Error fetching specializations", "err", err)
		return err
	}
	return Render(w, r, settings.ResultsSpecs(specializations))
}

func HandleSpecializationsEntry(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HandleSpecializationsEntry", "id_param", id_param)

	specialization, err := server.MyS.DB.SpecializationById(id_param)
	if err != nil {
		slog.Error("Error fetching specialization", "err", err)
		return err
	}

	return Render(w, r, settings.Entry_Spec(specialization))
}

func HandleSpecializationsDetails(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HSDetails", "id_param", id_param)
	spec_data, err := server.MyS.DB.SpecializationById(id_param)
	slog.Info("quere", "q", r.URL.Query())
	if err != nil {
		return err
	}
	slog.Info("HSettingsSpecDetails", "spec", spec_data)
	return Render(w, r, settings.Details_Spec(spec_data, types.SpecializationErrors{}))
}

func extractSpecsSortsFromForm(r *http.Request) *SpecFilter {
	return &SpecFilter{
		Search: r.FormValue("Search"),
		Order:  r.FormValue("Order"),
		SortBy: r.FormValue("SortBy"),
	}
}

func HandleSortingSpecs(w http.ResponseWriter, r *http.Request) error {
	log.Printf("All query parameters: %v", r.URL.Query()) //log do wywalenia

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return err
	}

	filter := *extractSpecsSortsFromForm(r)

	sortBy := filter.SortBy
	order := filter.Order
	search := filter.Search

	specs, err := server.MyS.DB.GetSortedSpecs(sortBy, order, search)
	if err != nil {
		http.Error(w, "Failed to fetch courses", http.StatusInternalServerError)
		return err
	}

	return Render(w, r, settings.EntrySpecsOnly(specs))
}

func HandleSpecializationsNew(w http.ResponseWriter, r *http.Request) error {
	spec := *extractSpecFromForm(r)

	errors, ok := validators.ValidateSpecialization(spec)
	if !ok {
		errors.Correct = false
		return Render(w, r, settings.NewEntrySwap_Spec(types.Specialization{}, spec, errors))
	}

	slog.Info("add spec", "correct", true)

	specId, err := server.MyS.DB.InsertSpecialization(spec)
	if err != nil {
		slog.Error("spec to db", "err", err)
		errors.Correct = false
		return Render(w, r, settings.NewEntrySwap_Spec(types.Specialization{}, spec, errors))
	}

	slog.Info("spec to db", "new_id", specId)
	spec.Id = int(specId)
	errors.Correct = true

	return Render(w, r, settings.NewEntrySwap_Spec(spec, types.Specialization{}, errors))
}

func HandleSpecializationsUpdate(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HandleSpecializationsUpdate", "id_param", id_param)

	specialization := types.Specialization{
		Name: r.FormValue("name"),
	}

	errors, ok := validators.ValidateSpecialization(specialization)
	if !ok {
		errors.Correct = false
		return Render(w, r, settings.Entry_Spec(specialization))
	}

	var err error
	specialization.Id, err = strconv.Atoi(id_param)
	if err != nil {
		slog.Error("Error parsing specialization ID", "err", err)
		return Render(w, r, settings.Entry_Spec(specialization))
	}

	err = server.MyS.DB.UpdateSpecialization(specialization)
	if err != nil {
		slog.Error("Error updating specialization", "err", err)
		return Render(w, r, settings.Entry_Spec(specialization))
	}

	slog.Info("Specialization updated", "specializationId", specialization.Id)
	errors.Correct = true
	return Render(w, r, settings.Entry_Spec(specialization))
}

func extractSpecFromForm(r *http.Request) *types.Specialization {
	return &types.Specialization{
		Name: r.FormValue("name"),
	}
}

func HandleSpecializationsGetNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, settings.NewEntry_Spec(types.Specialization{}, types.SpecializationErrors{}))
}

func HandleSpecializationsClearNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, settings.EmptySpace_Spec())
}
