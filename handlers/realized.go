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
	"thesis-management-app/views/realized"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/go-chi/chi/v5"
)

func HandleRealized(w http.ResponseWriter, r *http.Request) error {
	thes_data, err := server.MyS.DB.AllRealizedThesisEntries("thesis_id", true, r.URL.Query())
	if err != nil {
		slog.Error("HandleRealized", "err", err)
		return err
	}
	return Render(w, r, realized.Index(thes_data))
}

func HandleRealizedGenerateExcel(w http.ResponseWriter, r *http.Request) error {
	queryParams := r.URL.Query()
	filePath := "/realized/generate_excel"
	redirectURL := filePath + "?" + queryParams.Encode()

	currentTime := time.Now()
	fileName := "Wybrane Prace " + currentTime.Format("2-01-2006 15h04m05s") + ".xlsx"

	println(fileName)

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("HX-Redirect", redirectURL)

	t_data, err := filterRealizedThesisEntries(r)
	if err != nil {
		slog.Error("HandleRealizedGenerateExcel", "err", err)
		return err
	}

	f := excelize.NewFile()

	sheetName := "Prace"
	sheetIndex, _ := f.NewSheet(sheetName)

	headers := []string{
		"Numer Pracy", "Data Egzaminu", "Średnia Ocen ze Studiów", "Ocena z Egzaminu Kompetecyjnego",
		"Ocena z Egzaminu dyplomowego", "Ostateczny Wynik Studiów", "Ostateczny Wynik Studiów (słownie)",
		"Tytuł Pracy Dyplomowej (polski)", "Tytuł Pracy Dyplomowej (angielski)",
		"Język Pracy", "Biblioteka", "Numer Indeksu", "Imię Studenta", "Nazwisko Studenta",
		"Kierunek", "Specjalność", "Tryb", "Tytuł Naukowy", "Przewodniczący Imię", "Przewodniczący Nazwisko",
		"Tytuł Naukowy", "Promotor Imię", "Promotor Nazwisko",
		"Tytuł Naukowy", "Promotor Pomocniczy Imię", "Promotor Pomocniczy Nazwisko",
		"Tytuł Naukowy", "Recenzent Imię", "Recenzent Nazwisko",
		"Recenzent Godziny Rozliczeń", "Promotor Godziny Rozliczeń", "Promotor Pomocniczy Godziny Rozliczeń"}

	for i, header := range headers {
		col := ""
		index := i
		for index >= 0 {
			col = string(rune('A'+index%26)) + col
			index = index/26 - 1
		}
		err := f.SetCellValue(sheetName, col+"1", header)
		if err != nil {
			slog.Error("HandleRealizedGenerateExcel", "err", err)
			return err
		}
	}

	for i, t := range t_data {
		row := strconv.Itoa(i + 2) // Starts from second row
		data := map[string]interface{}{
			"A":  t.ThesisNumber,
			"B":  t.ExamDate,
			"C":  t.AverageStudyGrade,
			"D":  t.CompetencyExamGrade,
			"E":  t.DiplomaExamGrade,
			"F":  t.FinalStudyResult,
			"G":  t.FinalStudyResultText,
			"H":  t.ThesisTitlePolish,
			"I":  t.ThesisTitleEnglish,
			"J":  t.ThesisLanguage,
			"K":  t.Library,
			"L":  t.Student.StudentNumber,
			"M":  t.Student.FirstName,
			"N":  t.Student.LastName,
			"O":  t.Student.FieldOfStudy,
			"P":  t.Student.Specialization,
			"Q":  t.Student.ModeOfStudies,
			"R":  t.Chair.CurrentAcademicTitle,
			"S":  t.Chair.FirstName,
			"T":  t.Chair.LastName,
			"U":  t.Supervisor.CurrentAcademicTitle,
			"V":  t.Supervisor.FirstName,
			"W":  t.Supervisor.LastName,
			"X":  t.AssistantSupervisor.CurrentAcademicTitle,
			"Y":  t.AssistantSupervisor.FirstName,
			"Z":  t.AssistantSupervisor.LastName,
			"AA": t.Reviewer.CurrentAcademicTitle,
			"AB": t.Reviewer.FirstName,
			"AC": t.Reviewer.LastName,
			"AD": t.HourlySettlement.ReviewerHours,
			"AE": t.HourlySettlement.SupervisorHours,
			"AF": t.HourlySettlement.AssistantSupervisorHours,
		}

		for col, value := range data {
			cell := col + row
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				return fmt.Errorf("error setting cell %s: %w", cell, err)
			}
		}
	}

	f.SetActiveSheet(sheetIndex)

	if err := f.Write(w); err != nil {
		slog.Error("RealizedThesisGenerateExcel", "err", err)
	}
	return nil
}

func HandleAutocompleteThesisTitlePolish(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("thesis_title")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllThesisTitlesPolish(userInput)
	if err != nil {
		slog.Error("HandleAutocompleteThesisTitlePolish", "err", err)
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteStudentSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("student_name")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllStudentSurnames(userInput)
	if err != nil {
		slog.Error("HandleAutocompleteStudentSurname", "err", err)
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}
func HandleAutocompleteStudentNumber(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("student_number")

	if userInput == "" {
		return nil
	}

	queryResult, err := server.MyS.DB.GetAllStudentNumbers(userInput)
	if err != nil {
		slog.Error("HandleAutocompleteStudentNumber", "err", err)
		return err
	}

	for _, title := range queryResult {
		fmt.Println("Thesis Title:", title)
	}

	maxResults := 6
	if len(queryResult) > maxResults {
		queryResult = queryResult[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range queryResult {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteStudentNameAndSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("StudentNameAndSurname")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllStudentsNamesAndSurnames(userInput)
	if err != nil {
		return err
	}
	for _, title := range filteredThesisTitlesPolish {
		fmt.Println(title)
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteSupervisorName(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("firstNameSupervisor")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesNames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}



func HandleAutocompleteSupervisorSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("lastNameSupervisor")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteSupervisorTitle(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("supervisorAcademicTitle")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesTitles(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}
func HandleAutocompleteAssistantSupervisorName(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("firstNameAssistantSupervisor")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesNames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}




func HandleAutocompleteAssistantSupervisorSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("lastNameAssistantSupervisor")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteSupervisorNameAndSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("SupervisorNameAndSurname")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}
	for _, title := range filteredThesisTitlesPolish {
		fmt.Println(title)
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteAssistantSupervisorTitle(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("assistantSupervisorAcademicTitle")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesTitles(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}





func HandleAutocompleteReviewerName(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("firstNameReviewer")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesNames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}




func HandleAutocompleteReviewerSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("lastNameReviewer")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteReviewerTitle(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("reviewerAcademicTitle")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesTitles(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}



func HandleAutocompleteChairName(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("firstNameChair")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesNames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}




func HandleAutocompleteChairSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("lastNameChair")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteChairTitle(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("chairAcademicTitle")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllUniversityEmployeesTitles(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}
func HandleAutocompleteCourse(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("course")

	if userInput == "" {
		return nil
	}

	filteredThesisTitlesPolish, err := server.MyS.DB.GetAllCourseNames(userInput)
	if err != nil {
		slog.Error("HandleAutocompleteCourse", "err", err)
		return err
	}
	for _, name := range filteredThesisTitlesPolish {
		fmt.Println(name)
	}

	maxResults := 6
	if len(filteredThesisTitlesPolish) > maxResults {
		filteredThesisTitlesPolish = filteredThesisTitlesPolish[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredThesisTitlesPolish {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func filterRealizedThesisEntries(r *http.Request) ([]types.RealizedThesisEntry, error) {
	q := r.URL.Query()
	sortBy := "thesis_id"
	desc := true
	searchString := ""
	for key, val := range q {
		if val[0] == "" {
			q.Del(key)
		}
		if key == "SortBy" {
			sortBy = val[0]
			q.Del(key)
		}
		if key == "Order" {
			if val[0] == "ASC" {
				desc = false
			}
			q.Del(key)
		}
		if key == "Search" {
			searchString = val[0]
			q.Del(key)
		}
	}
	slog.Info("filterRealizedThesisEntries", "searchString", searchString)
	r.URL.RawQuery = q.Encode()
	thes_data, err := server.MyS.DB.AllRealizedThesisEntries(sortBy, desc, r.URL.Query())
	if err != nil {
		return nil, err
	}
	if searchString == "" {
		return thes_data, err
	}
	results := []types.RealizedThesisEntry{}
	for _, t := range thes_data {
		lookupString := strings.ToLower(fmt.Sprintf("%v", t))
		slog.Info("filterRealizedThesisEntries", "lookupString", lookupString)
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
	return results, nil
}

func HandleRealizedFiltered(w http.ResponseWriter, r *http.Request) error {
	results, err := filterRealizedThesisEntries(r)
	if err != nil {
		return err
	}
	return Render(w, r, realized.Results(results))
}

func HandleRealizedDetails(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HRDetails", "id_param", id_param)
	thes_data, err := server.MyS.DB.RealizedThesisEntryByID(id_param)
	slog.Info("HRDetails", "q", r.URL.Query())
	slog.Info("HRDetails", "thes", thes_data)
	if err != nil {
		slog.Error("HRDetails", "err", err)
		return err
	}
	return Render(w, r, realized.Details(thes_data, types.RealizedThesisEntryErrors{}))
}

func HandleRealizedEntry(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HREntry", "id_param", id_param)
	thes_data, err := server.MyS.DB.RealizedThesisEntryByID(id_param)
	if err != nil {
		return err
	}
	return Render(w, r, realized.Entry(thes_data))
}

func extractRealizedThesisFromForm(r *http.Request) *types.RealizedThesisEntry {
	return &types.RealizedThesisEntry{
		ThesisNumber:         r.FormValue("thesisNumber"),
		ExamDate:             r.FormValue("examDate"),
		AverageStudyGrade:    r.FormValue("averageStudyGrade"),
		CompetencyExamGrade:  r.FormValue("competencyExamGrade"),
		DiplomaExamGrade:     r.FormValue("diplomaExamGrade"),
		FinalStudyResult:     r.FormValue("finalStudyResult"),
		FinalStudyResultText: r.FormValue("finalStudyResultText"),
		ThesisTitlePolish:    r.FormValue("thesisTitlePolish"),
		ThesisTitleEnglish:   r.FormValue("thesisTitleEnglish"),
		ThesisLanguage:       r.FormValue("thesisLanguage"),
		Library:              r.FormValue("library"),
		Student: types.Student{
			StudentNumber:  r.FormValue("studentNumber"),
			FirstName:      r.FormValue("firstNameStudent"),
			LastName:       r.FormValue("lastNameStudent"),
			FieldOfStudy:   r.FormValue("fieldOfStudy"),
			Specialization: r.FormValue("specialization"),
			ModeOfStudies:  r.FormValue("modeOfStudies"),
		},
		ChairAcademicTitle: r.FormValue("chairAcademicTitle"),
		Chair: types.UniversityEmployee{
			FirstName:            r.FormValue("firstNameChair"),
			LastName:             r.FormValue("lastNameChair"),
			CurrentAcademicTitle: r.FormValue("chairAcademicTitle"),
		},
		SupervisorAcademicTitle: r.FormValue("supervisorAcademicTitle"),
		Supervisor: types.UniversityEmployee{
			FirstName:            r.FormValue("firstNameSupervisor"),
			LastName:             r.FormValue("lastNameSupervisor"),
			CurrentAcademicTitle: r.FormValue("supervisorAcademicTitle"),
		},
		AssistantSupervisorAcademicTitle: r.FormValue("assistantSupervisorAcademicTitle"),
		AssistantSupervisor: types.UniversityEmployee{
			FirstName:            r.FormValue("firstNameAssistantSupervisor"),
			LastName:             r.FormValue("lastNameAssistantSupervisor"),
			CurrentAcademicTitle: r.FormValue("assistantSupervisorAcademicTitle"),
		},
		ReviewerAcademicTitle: r.FormValue("reviewerAcademicTitle"),
		Reviewer: types.UniversityEmployee{
			FirstName:            r.FormValue("firstNameReviewer"),
			LastName:             r.FormValue("lastNameReviewer"),
			CurrentAcademicTitle: r.FormValue("reviewerAcademicTitle"),
		},
		HourlySettlement: types.HourlySettlement{},
	}
}

func getEmployeeId(emp types.UniversityEmployee) (int, error) {
	empId, err := server.MyS.DB.EmployeeIdByName(emp.FirstName + " " + emp.LastName)
	slog.Info("getEmployeeId", "empName", emp.FirstName+" "+emp.LastName)
	slog.Info("getEmployeeId", "empId", empId)
	if err != nil {
		slog.Error("getEmployeeId", "err", err)
	}
	if empId == 0 {
		if emp.FirstName != "" && emp.LastName != "" {
			var id int64
			id, err = server.MyS.DB.InsertUniversityEmployee(emp)
			if err != nil {
				slog.Error("getEmployeeId", "err", err)
			}
			slog.Info("getEmployeeId", "inserting new emp, id", id)
			empId = int(id)
		}
	}
	return empId, err
}

func HandleRealizedNew(w http.ResponseWriter, r *http.Request) error {
	t := *extractRealizedThesisFromForm(r)
	errors, ok := validators.ValidateRealizedThesis(t)
	if !ok {
		errors.Correct = false
		return Render(w, r, realized.NewEntrySwap(types.RealizedThesisEntry{}, t, errors))
	}
	sId, err := server.MyS.DB.InsertStudent(t.Student)
	if err != nil {
		slog.Error("student to db", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.NewEntrySwap(types.RealizedThesisEntry{}, t, errors))
	}
	t.Student.Id = int(sId)
	supId, err := getEmployeeId(t.Supervisor)
	if err != nil {
		slog.Error("InsertNew", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.NewEntrySwap(types.RealizedThesisEntry{}, t, errors))
	}
	t.Supervisor.Id = supId
	asId, err := getEmployeeId(t.AssistantSupervisor)
	if err != nil {
		slog.Error("InsertNew", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.NewEntrySwap(types.RealizedThesisEntry{}, t, errors))
	}
	t.AssistantSupervisor.Id = asId
	reId, err := getEmployeeId(t.Reviewer)
	if err != nil {
		slog.Error("InsertNew", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.NewEntrySwap(types.RealizedThesisEntry{}, t, errors))
	}
	t.Reviewer.Id = reId
	chId, err := getEmployeeId(t.Chair)
	if err != nil {
		slog.Error("InsertNew", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.NewEntrySwap(types.RealizedThesisEntry{}, t, errors))
	}
	t.Chair.Id = chId
	tId, err := server.MyS.DB.InsertRealizedThesisByEntry(&t)
	if err != nil {
		slog.Error("InsertNew", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.NewEntrySwap(types.RealizedThesisEntry{}, t, errors))
	}
	slog.Info("thesis to db", "new_id", tId)
	slog.Info("sudent to db", "new_id", sId)
	errors.Correct = true
	return Render(w, r, realized.NewEntrySwap(t, types.RealizedThesisEntry{}, errors))
}

func HandleRealizedUpdate(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("UPDATE", "id_param", id_param)
	t := *extractRealizedThesisFromForm(r)
	var err error
	t.Id, err = strconv.Atoi(id_param)
	slog.Info("UpdateRealizedThesis", "t before", t)
	errors, ok := validators.ValidateRealizedThesis(t)
	if !ok {
		errors.Correct = false
		return Render(w, r, realized.Details(t, errors))
	}
	if err != nil {
		slog.Error("Update", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.Details(t, errors))
	}
	t.Student.Id, err = server.MyS.DB.GetStudentIdFromThesisEntry(t.Id)
	slog.Info("UpdateRealizedThesis", "student_id", t.Student.Id)
	if err != nil {
		slog.Error("Update get stud id", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.Details(t, errors))
	}
	sId, err := server.MyS.DB.UpdateStudent(t.Student)
	if err != nil {
		slog.Error("Update stud", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.Details(t, errors))
	}
	t.Student.Id = int(sId)
	supId, err := getEmployeeId(t.Supervisor)
	slog.Info("UpdateRealizedThesis", "supervisor", t.Supervisor)
	slog.Info("UpdateRealizedThesis", "supervisor_id", supId)
	if err != nil {
		slog.Error("Update emp", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.Details(t, errors))
	}
	t.Supervisor.Id = supId
	asId, err := getEmployeeId(t.AssistantSupervisor)
	slog.Info("UpdateRealizedThesis", "assistant_supervisor", t.AssistantSupervisor)
	slog.Info("UpdateRealizedThesis", "assistant_supervisor_id", asId)
	if err != nil {
		slog.Error("Update emp", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.Details(t, errors))
	}
	t.AssistantSupervisor.Id = asId
	reId, err := getEmployeeId(t.Reviewer)
	slog.Info("UpdateRealizedThesis", "Reviewer_id", reId)
	if err != nil {
		slog.Error("update emp", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.Details(t, errors))
	}
	t.Reviewer.Id = reId
	chId, err := getEmployeeId(t.Chair)
	if err != nil {
		slog.Error("update emp", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.Details(t, errors))
	}
	t.Chair.Id = chId
	err = server.MyS.DB.UpdateRealizedThesisByEntry(&t)
	if err != nil {
		slog.Error("Update Thesis", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.Details(t, errors))
	}
	slog.Info("update thesis", "correct", true)
	errors.Correct = true
	slog.Info("UpdateRealizedThesis", "t after", t)
	return Render(w, r, realized.Entry(t))
}

func HandleRealizedGetNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, realized.NewEntry(types.RealizedThesisEntry{}, types.RealizedThesisEntryErrors{}))
}

func HandleRealizedClearNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, realized.EmptySpace())
}
