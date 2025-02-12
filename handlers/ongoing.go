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
	"time"

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
	// slog.Info("HandleOngoingDetails", "id_param", id_param)
	thes_data, err := server.MyS.DB.OngoingThesisEntryByID(id_param)
	// slog.Info("HandleOngoingDetails", "q", r.URL.Query())
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
	//t.ThesisNumber = validators.CheckThesisNumber(t.ThesisNumber, t.Student.Degree)
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
	_, err = server.MyS.DB.InsertTask(types.Task{Content: "Podanie o egzamin dyplomowy", OngoingThesisID: t.Id})
	_, err = server.MyS.DB.InsertTask(types.Task{Content: "Oświadczenie o badaniach losów zawodowych absolwenta", OngoingThesisID: t.Id})
	_, err = server.MyS.DB.InsertTask(types.Task{Content: "Obiegówka", OngoingThesisID: t.Id})
	_, err = server.MyS.DB.InsertTask(types.Task{Content: "Wydruk prac dyplomowych", OngoingThesisID: t.Id})
	if t.Student.Degree == "II stopień" {
		_, err = server.MyS.DB.InsertTask(types.Task{Content: "Oświadczenie o nieużywaniu legitymacji", OngoingThesisID: t.Id})
	}
	if t.Student.ModeOfStudies == "niestacjonarne" {
		_, err = server.MyS.DB.InsertTask(types.Task{Content: "Podanie do prodziekana o wyznaczenie egzaminu w dzień roboczy", OngoingThesisID: t.Id})
	}
	return Render(w, r, ongoing.NewEntrySwap(t, types.OngoingThesisEntry{}, errors))
}

func HandleOngoingFiltered(w http.ResponseWriter, r *http.Request) error {
	// slog.Info("HandleOngoingFiltered", "query", r.URL.Query())
	q := r.URL.Query()
	page, err := strconv.Atoi(q.Get("page_number"))
	if err != nil {
		slog.Info("HandleOngoingFiltered", "err", err)
		return err
	}
	pageSize, err := strconv.Atoi(q.Get("page_size"))
	if err != nil {
		slog.Info("HandleOngoingFiltered", "err", err)
		return err
	}
	if q.Get("reset_page") != "false" {
		q.Set("page_number", strconv.Itoa(1))
		page = 1
	}
	q.Del("reset_page")
	r.URL.RawQuery = q.Encode()
	results, err := filterOngoingThesisEntries(r)
	if err != nil {
		return err
	}
	return Render(w, r, ongoing.SwapResults(results, page, pageSize))
}

func HandleOngoingUpdate(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("UPDATE", "id_param", id_param)
	user, ok := r.Context().Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		slog.Error("UpdateOngoing", "Failed To retrieve user", user)
	}
	t := *extractOngoingThesisFromForm(r)
	var err error
	t.Id, err = strconv.Atoi(id_param)
	errors, ok := validators.ValidateOngoingThesis(t)
	if user.IsAdmin {
		if !ok {
			errors.Correct = false
			return Render(w, r, ongoing.Details(t, errors))
		}
		if err != nil {
			slog.Error("Update", "err", err)
			errors.InternalError = true
			return Render(w, r, ongoing.Details(t, errors))
		}
		t.Student.Id, err = server.MyS.DB.GetStudentIdFromOngoingThesisEntry(t.Id)
		// slog.Info("UpdateOngoingThesis", "student_id", t.Student.Id)
		if err != nil {
			slog.Error("Update get stud id", "err", err)
			errors.InternalError = true
			return Render(w, r, ongoing.Details(t, errors))
		}
		sId, err := server.MyS.DB.UpdateStudent(t.Student)
		if err != nil {
			slog.Error("Update stud", "err", err)
			errors.InternalError = true
			return Render(w, r, ongoing.Details(t, errors))
		}
		t.Student.Id = int(sId)
		supId, err := getEmployeeId(t.Supervisor)
		if err != nil {
			slog.Error("Update emp", "err", err)
			errors.InternalError = true
			return Render(w, r, ongoing.Details(t, errors))
		}
		t.Supervisor.Id = supId
		asId, err := getEmployeeId(t.AssistantSupervisor)
		if err != nil {
			slog.Error("Update emp", "err", err)
			errors.InternalError = true
			return Render(w, r, ongoing.Details(t, errors))
		}
		t.AssistantSupervisor.Id = asId
		err = server.MyS.DB.UpdateOngoingThesisByEntry(&t)
		if err != nil {
			slog.Error("Update Thesis", "err", err)
			errors.InternalError = true
			return Render(w, r, ongoing.Details(t, errors))
		}
	}
	err = updateNoteOngoing(t.Note.Content, t.Id, user.Id)
	if err != nil {
		slog.Error("update note", "err", err)
		errors.InternalError = true
		return Render(w, r, ongoing.Details(t, errors))
	}
	errors.Correct = true
	return Render(w, r, ongoing.Entry(t))
}

func HandleOngoingArchive(w http.ResponseWriter, r *http.Request) error {
	thesisId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	thesis, err := server.MyS.DB.OngoingThesisEntryByID(strconv.Itoa(thesisId))
	if err != nil {
		return err
	}
	//TODO: check if all tasks are completed
	if !server.MyS.DB.CheckIfAllTaskAreCompleted(thesisId) {
		return Render(w, r, ongoing.Details(thesis, types.OngoingThesisEntryErrors{Checklist: "Wszystkie podpunkty muszą być zaznaczone aby zarchiwizować"}))
	}
	err = server.MyS.DB.ArchiveOngoingThesis(thesisId)
	if err != nil {
		return err
	}
	thesis.Archived = "true"
	// add thesis to realized
	thesis.ThesisNumber = validators.CheckThesisNumber(fmt.Sprintf("k22/stopien/num/%v", time.Now().Year()), thesis.Student.Degree)
	newId, err := server.MyS.DB.InsertOngoingThesisToRealized(thesis)
	if err != nil {
		return err
	}
	thesis.Note.RealizedThesisID = int(newId)
	thesis.Note.OngoingThesisID = thesisId
	// slog.Info("ArchiveOngoingThesis", "note", thesis.Note)
	err = server.MyS.DB.UpdateNoteRelatedOngoingThesis(thesis.Note)
	if err != nil {
		return err
	}
	return Render(w, r, ongoing.Entry(thesis))
}

func HandleOngoingNext(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	num := query.Get("page_number")
	// slog.Info("HandleOngoingNext", "page", num)
	// slog.Info("HandleOngoingNext", "query", r.URL.Query())
	page, err := strconv.Atoi(num)
	if err != nil {
		slog.Error("HandleOngoingNext", "err", err)
	}
	query.Set("page_number", strconv.Itoa(page+1))
	query.Set("reset_page", "false")
	r.URL.RawQuery = query.Encode()
	return HandleOngoingFiltered(w, r)
}

func HandleOngoingPrev(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	num := query.Get("page_number")
	// slog.Info("HandleOngoingNext", "page", num)
	// slog.Info("HandleOngoingNext", "query", r.URL.Query())
	page, err := strconv.Atoi(num)
	if err != nil {
		slog.Error("HandleOngoingNext", "err", err)
	}
	query.Set("page_number", strconv.Itoa(page-1))
	query.Set("reset_page", "false")
	r.URL.RawQuery = query.Encode()
	return HandleOngoingFiltered(w, r)
}

func HandleOngoingGetNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, ongoing.NewEntry(types.OngoingThesisEntry{}, types.OngoingThesisEntryErrors{}))
}

func HandleOngoingClearNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, ongoing.EmptySpace())
}

func extractOngoingThesisFromForm(r *http.Request) *types.OngoingThesisEntry {
	return &types.OngoingThesisEntry{
		//ThesisNumber:       strings.TrimSpace(r.FormValue("thesisNumber")),
		ThesisTitlePolish:  strings.TrimSpace(r.FormValue("thesisTitlePolish")),
		ThesisTitleEnglish: strings.TrimSpace(r.FormValue("thesisTitleEnglish")),
		ThesisLanguage:     strings.TrimSpace(r.FormValue("thesisLanguage")),
		Student: types.Student{
			StudentNumber:  strings.TrimSpace(r.FormValue("studentNumber")),
			FirstName:      strings.TrimSpace(r.FormValue("firstNameStudent")),
			LastName:       strings.TrimSpace(r.FormValue("lastNameStudent")),
			FieldOfStudy:   strings.TrimSpace(r.FormValue("course")),
			Specialization: strings.TrimSpace(r.FormValue("specialization")),
			ModeOfStudies:  strings.TrimSpace(r.FormValue("modeOfStudies")),
			Degree:         strings.TrimSpace(r.FormValue("degree")),
		},
		SupervisorAcademicTitle: strings.TrimSpace(r.FormValue("supervisorAcademicTitle")),
		Supervisor: types.UniversityEmployeeEntry{
			FirstName:            strings.TrimSpace(r.FormValue("firstNameSupervisor")),
			LastName:             strings.TrimSpace(r.FormValue("lastNameSupervisor")),
			CurrentAcademicTitle: strings.TrimSpace(r.FormValue("supervisorAcademicTitle")),
		},
		AssistantSupervisorAcademicTitle: strings.TrimSpace(r.FormValue("assistantSupervisorAcademicTitle")),
		AssistantSupervisor: types.UniversityEmployeeEntry{
			FirstName:            strings.TrimSpace(r.FormValue("firstNameAssistantSupervisor")),
			LastName:             strings.TrimSpace(r.FormValue("lastNameAssistantSupervisor")),
			CurrentAcademicTitle: strings.TrimSpace(r.FormValue("assistantSupervisorAcademicTitle")),
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
	page_size := 20
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
		} else if key == "page_size" {
			page_size, _ = strconv.Atoi(val[0])
			slog.Info("filterOngoingThesisEntries", "page_size", page_size)
			q.Del(key)
		} else if key == "Search" {
			searchString = val[0]
			q.Del(key)
		}
	}
	slog.Info("filterOngoingThesisEntries", "searchString", searchString)
	r.URL.RawQuery = q.Encode()
	if searchString == "" {
		thes_data, err := server.MyS.DB.AllOngoingThesisEntries(sortBy, desc, page_num, page_size, r.URL.Query())
		if err != nil {
			return nil, err
		}
		return thes_data, nil
	}
	thes_data, err := server.MyS.DB.AllOngoingThesisEntries(sortBy, desc, -1, page_size, r.URL.Query())
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
	paginated_res, _ := paginate(results, page_num, page_size)
	return paginated_res, nil
}

func updateNoteOngoing(newContent string, thesisId int, userId int) error {
	note, err := server.MyS.DB.GetNote(0, thesisId, userId)
	if err != nil {
		slog.Error("updateNoteOngoing", "err", err)
		return fmt.Errorf("updateNoteOngoing occured error while retrieving note: %v", err)
	}
	if note.Id == 0 {
		_, err := server.MyS.DB.InsertNote(types.Note{Content: newContent, OngoingThesisID: thesisId, UniversityEmployeeID: userId})
		if err != nil {
			slog.Error("updateNoteOngoing", "err", err)
			return fmt.Errorf("updateNoteOngoing occured error while inserting note: %v", err)
		}
		return nil
	}
	note.Content = newContent
	err = server.MyS.DB.UpdateNote(note)
	if err != nil {
		slog.Error("updateNoteOngoing", "err", err)
		return fmt.Errorf("updateNoteOngoing occured error while updating note: %v", err)
	}
	return nil
}
