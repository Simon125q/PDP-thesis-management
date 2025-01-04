package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/pkgs/validators"
	"thesis-management-app/types"
	"thesis-management-app/views/ongoing"

	"github.com/go-chi/chi/v5"
)

func HandleOngoing(w http.ResponseWriter, r *http.Request) error {
	slog.Info("HandleOngoing", "entered", true)
	thes_data, err := server.MyS.DB.AllOngoingThesisEntries("thesis_id", true, 1, PageLimit, r.URL.Query())
	if err != nil {
		slog.Error("HandleOngoing", "err", err)
		return err
	}
	return Render(w, r, ongoing.Index(thes_data))
}

func HandleOngoingDetails(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HandleOngoingDetails", "id_param", id_param)
	thes_data, err := server.MyS.DB.OngoingThesisEntryByID(id_param)
	slog.Info("HandleOngoingDetails", "q", r.URL.Query())
	if err != nil {
		slog.Error("HandleOngoingDetails", "err", err)
		return fmt.Errorf("HandleOngoingDetails -> %v", err)
	}
	return Render(w, r, ongoing.Details(thes_data, types.OngoingThesisEntryErrors{}))
}

func HandleOngoingEntry(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HandleOngoingEntry", "id_param", id_param)
	thes_data, err := server.MyS.DB.OngoingThesisEntryByID(id_param)
	if err != nil {
		return err
	}
	return Render(w, r, ongoing.Entry(thes_data))
}

func HandleOngoingNew(w http.ResponseWriter, r *http.Request) error {
	t := *extractOngoingThesisFromForm(r)
	t.ThesisNumber = validators.CheckThesisNumber(t.ThesisNumber, t.Student.Degree)
	errors, ok := validators.ValidateOngoingThesis(t)
	if !ok {
		errors.Correct = false
		return Render(w, r, ongoing.NewEntrySwap(types.OngoingThesisEntry{}, t, errors))
	}
	sId, err := server.MyS.DB.InsertStudent(t.Student)
	if err != nil {
		slog.Error("student to db", "err", err)
		errors.InternalError = true
		return Render(w, r, ongoing.NewEntrySwap(types.OngoingThesisEntry{}, t, errors))
	}
	t.Student.Id = int(sId)
	supId, err := getEmployeeId(t.Supervisor)
	if err != nil {
		slog.Error("Insert Employee New Ongoing", "err", err)
		errors.InternalError = true
		return Render(w, r, ongoing.NewEntrySwap(types.OngoingThesisEntry{}, t, errors))
	}
	t.Supervisor.Id = supId
	asId, err := getEmployeeId(t.AssistantSupervisor)
	if err != nil {
		slog.Error("Insert Employee New Ongoing", "err", err)
		errors.InternalError = true
		return Render(w, r, ongoing.NewEntrySwap(types.OngoingThesisEntry{}, t, errors))
	}
	t.AssistantSupervisor.Id = asId
	tId, err := server.MyS.DB.InsertOngoingThesisByEntry(&t)
	if err != nil {
		slog.Error("Insert New Ongoing Thesis", "err", err)
		errors.InternalError = true
		return Render(w, r, ongoing.NewEntrySwap(types.OngoingThesisEntry{}, t, errors))
	}
	t.Id = int(tId)
	slog.Info("ongoing thesis to db", "new_id", tId)
	errors.Correct = true
	return Render(w, r, ongoing.NewEntrySwap(t, types.OngoingThesisEntry{}, errors))
}

func HandleOngoingFiltered(w http.ResponseWriter, r *http.Request) error {
	slog.Info("HandleOngoingFiltered", "query", r.URL.Query())
	q := r.URL.Query()
	page, err := strconv.Atoi(q.Get("page_number"))
	if q.Get("reset_page") != "false" {
		q.Set("page_number", strconv.Itoa(1))
		page = 1
	}
	q.Del("reset_page")
	r.URL.RawQuery = q.Encode()
	if err != nil {
		slog.Info("HandleOngoingFiltered", "err", err)
		return err
	}
	results, err := filterOngoingThesisEntries(r)
	if err != nil {
		return err
	}
	return Render(w, r, ongoing.SwapResults(results, page))
}

func HandleOngoingGetNew(w http.ResponseWriter, r *http.Request) error {
	slog.Info("HandleOngoingGetNew", "entered", true)
	return Render(w, r, ongoing.NewEntry(types.OngoingThesisEntry{}, types.OngoingThesisEntryErrors{}))
}

func HandleOngoingClearNew(w http.ResponseWriter, r *http.Request) error {
	slog.Info("HandleOngoingClearNew", "entered", true)
	return Render(w, r, ongoing.EmptySpace())
}

func extractOngoingThesisFromForm(r *http.Request) *types.OngoingThesisEntry {
	return &types.OngoingThesisEntry{
		ThesisNumber:       r.FormValue("thesisNumber"),
		ThesisTitlePolish:  r.FormValue("thesisTitlePolish"),
		ThesisTitleEnglish: r.FormValue("thesisTitleEnglish"),
		ThesisLanguage:     r.FormValue("thesisLanguage"),
		Student: types.Student{
			StudentNumber:  r.FormValue("studentNumber"),
			FirstName:      r.FormValue("firstNameStudent"),
			LastName:       r.FormValue("lastNameStudent"),
			FieldOfStudy:   r.FormValue("course"),
			Specialization: r.FormValue("specialization"),
			ModeOfStudies:  r.FormValue("modeOfStudies"),
			Degree:         r.FormValue("degree"),
		},
		SupervisorAcademicTitle: r.FormValue("supervisorAcademicTitle"),
		Supervisor: types.UniversityEmployeeEntry{
			FirstName:            r.FormValue("firstNameSupervisor"),
			LastName:             r.FormValue("lastNameSupervisor"),
			CurrentAcademicTitle: r.FormValue("supervisorAcademicTitle"),
		},
		AssistantSupervisorAcademicTitle: r.FormValue("assistantSupervisorAcademicTitle"),
		AssistantSupervisor: types.UniversityEmployeeEntry{
			FirstName:            r.FormValue("firstNameAssistantSupervisor"),
			LastName:             r.FormValue("lastNameAssistantSupervisor"),
			CurrentAcademicTitle: r.FormValue("assistantSupervisorAcademicTitle"),
		},
		Note: types.Note{
			Content: r.FormValue("thesis_note"),
		},
	}
}

func filterOngoingThesisEntries(r *http.Request) ([]types.OngoingThesisEntry, error) {
	q := r.URL.Query()
	sortBy := "thesis_id"
	desc := true
	searchString := ""
	page_num := 1
	for key, val := range q {
		if val[0] == "" || val[0] == "all" {
			q.Del(key)
		}
		if key == "SortBy" {
			sortBy = val[0]
			q.Del(key)
		} else if key == "Order" {
			if val[0] == "ASC" {
				desc = false
			}
			q.Del(key)
		} else if key == "page_number" {
			page_num, _ = strconv.Atoi(val[0])
			slog.Info("filterOngoingThesisEntries", "page_number", page_num)
			q.Del(key)
		} else if key == "Search" {
			searchString = val[0]
			q.Del(key)
		}
	}
	slog.Info("filterOngoingThesisEntries", "searchString", searchString)
	r.URL.RawQuery = q.Encode()
	if searchString == "" {
		thes_data, err := server.MyS.DB.AllOngoingThesisEntries(sortBy, desc, page_num, PageLimit, r.URL.Query())
		if err != nil {
			return nil, err
		}
		return thes_data, nil
	}
	thes_data, err := server.MyS.DB.AllOngoingThesisEntries(sortBy, desc, -1, PageLimit, r.URL.Query())
	if err != nil {
		return nil, err
	}
	results := []types.OngoingThesisEntry{}
	for _, t := range thes_data {
		lookupString := strings.ToLower(fmt.Sprintf("%v", t))
		slog.Info("filterOngoingThesisEntries", "lookupString", lookupString)
		match := true
		for _, part := range strings.Split(strings.ToLower(searchString), " ") {
			if !strings.Contains(lookupString, part) {
				match = false
				break
			}
		}
		if match {
			results = append(results, t)
		}
	}
	paginated_res, _ := paginate(results, page_num, PageLimit)
	return paginated_res, nil
}
