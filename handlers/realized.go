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
	"thesis-management-app/views/components"
	"thesis-management-app/views/realized"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/go-chi/chi/v5"
)

const PageLimit = 20

func HandleRealized(w http.ResponseWriter, r *http.Request) error {
	slog.Info("HandleRealized", "entered", true)
	thes_data, err := server.MyS.DB.AllRealizedThesisEntries("thesis_id", true, 1, PageLimit, r.URL.Query())
	if err != nil {
		slog.Error("HandleRealized", "err", err)
		return err
	}
	return Render(w, r, realized.Index(thes_data))
}

func HandleRealizedGenerateExcel(w http.ResponseWriter, r *http.Request) error {
	slog.Info("HandleRealizedGenerateExcel", "entered", true)
	queryParams := r.URL.Query()
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		fileName = "Wybrane_Prace"
	}

	queryParams.Set("page_number", "-1")
	r.URL.RawQuery = queryParams.Encode()
	filePath := "/realized/generate_excel"
	redirectURL := filePath + "?" + queryParams.Encode()
	queryParams.Del("fileName")
	r.URL.RawQuery = queryParams.Encode()
	if err := r.ParseForm(); err != nil {
		return err
	}

	currentTime := time.Now()
	fileName = fileName + "_" + currentTime.Format("2-01-2006_15h04m05s") + ".xlsx"

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("HX-Redirect", redirectURL)

	t_data, err := filterRealizedThesisEntries(r)
	if err != nil {
		slog.Error("HandleRealizedGenerateExcel", "err", err)
		return err
	}
	queryParams.Set("fileName", fileName)
	f := excelize.NewFile()

	sheetName := "Prace"
	if err := f.SetSheetName("Sheet1", sheetName); err != nil {
		slog.Error("Error changing sheet name", "err", err)
	}
	sheetIndex, _ := f.GetSheetIndex(sheetName)

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

	filteredResults, err := server.MyS.DB.GetAllThesisTitlesPolish(userInput)
	if err != nil {
		slog.Error("HandleAutocompleteThesisTitlePolish", "err", err)
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteStudentSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("student_name")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllStudentSurnames(userInput)
	if err != nil {
		slog.Error("HandleAutocompleteStudentSurname", "err", err)
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}
func HandleAutocompleteStudentNumber(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("student_number")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllStudentNumbers(userInput)
	if err != nil {
		slog.Error("HandleAutocompleteStudentNumber", "err", err)
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteStudentNameAndSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("student_name")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllStudentsNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteSupervisorNameAndSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("supervisor_name")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteAssistantSupervisorNameAndSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("assistant_supervisor_name")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteReviewerNameAndSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("reviewer_name")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteSupervisorName(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("firstNameSupervisor")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesTitlesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteSupervisorSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("lastNameSupervisor")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesTitlesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteSupervisorTitle(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("supervisorAcademicTitle")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesTitlesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}
func HandleAutocompleteAssistantSupervisorName(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("firstNameAssistantSupervisor")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesTitlesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteAssistantSupervisorSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("lastNameAssistantSupervisor")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesTitlesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteAssistantSupervisorTitle(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("assistantSupervisorAcademicTitle")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesTitlesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteReviewerName(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("firstNameReviewer")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesTitlesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteReviewerSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("lastNameReviewer")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesTitlesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteReviewerTitle(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("reviewerAcademicTitle")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesTitlesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteChairName(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("firstNameChair")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesTitlesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteChairSurname(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("lastNameChair")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesTitlesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteChairTitle(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("chairAcademicTitle")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllUniversityEmployeesTitlesNamesAndSurnames(userInput)
	if err != nil {
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}
func HandleAutocompleteCourse(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("course")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllCourseNames(userInput)
	if err != nil {
		slog.Error("HandleAutocompleteCourse", "err", err)
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func HandleAutocompleteSpecialization(w http.ResponseWriter, r *http.Request) error {

	userInput := r.URL.Query().Get("specialization")

	if userInput == "" {
		return nil
	}

	filteredResults, err := server.MyS.DB.GetAllSpecializationsNames(userInput)
	if err != nil {
		slog.Error("HandleAutocompleteCourse", "err", err)
		return err
	}

	maxResults := 6
	if len(filteredResults) > maxResults {
		filteredResults = filteredResults[:maxResults]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, title := range filteredResults {
		fmt.Fprintf(w, "<li class=\"suggestion px-3 py-2 hover:bg-gray-100 cursor-pointer w-full\" >%s</li>", title)
	}

	return nil
}

func filterRealizedThesisEntries(r *http.Request) ([]types.RealizedThesisEntry, error) {
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
			slog.Info("filterRealizedThesisEntries", "page_number", page_num)
			q.Del(key)
		} else if key == "Search" {
			searchString = val[0]
			q.Del(key)
		}
	}
	slog.Info("filterRealizedThesisEntries", "searchString", searchString)
	r.URL.RawQuery = q.Encode()
	if searchString == "" {
		thes_data, err := server.MyS.DB.AllRealizedThesisEntries(sortBy, desc, page_num, PageLimit, r.URL.Query())
		if err != nil {
			return nil, err
		}
		return thes_data, nil
	}
	thes_data, err := server.MyS.DB.AllRealizedThesisEntries(sortBy, desc, -1, PageLimit, r.URL.Query())
	if err != nil {
		return nil, err
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
	paginated_res, _ := paginate(results, page_num, PageLimit)
	return paginated_res, nil
}

func HandleRealizedNext(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	num := query.Get("page_number")
	slog.Info("HandleRealizedNext", "page", num)
	slog.Info("HandleRealizedNext", "query", r.URL.Query())
	page, err := strconv.Atoi(num)
	if err != nil {
		slog.Error("HandleRealizedNext", "err", err)
	}
	query.Set("page_number", strconv.Itoa(page+1))
	query.Set("reset_page", "false")
	r.URL.RawQuery = query.Encode()
	return HandleRealizedFiltered(w, r)
}

func HandleRealizedPrev(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	num := query.Get("page_number")
	slog.Info("HandleRealizedNext", "page", num)
	slog.Info("HandleRealizedNext", "query", r.URL.Query())
	page, err := strconv.Atoi(num)
	if err != nil {
		slog.Error("HandleRealizedNext", "err", err)
	}
	query.Set("page_number", strconv.Itoa(page-1))
	query.Set("reset_page", "false")
	r.URL.RawQuery = query.Encode()
	return HandleRealizedFiltered(w, r)
}

func HandleRealizedFiltered(w http.ResponseWriter, r *http.Request) error {
	slog.Info("HandleRealizedFiltered", "query", r.URL.Query())
	q := r.URL.Query()
	page, err := strconv.Atoi(q.Get("page_number"))
	if q.Get("reset_page") != "false" {
		q.Set("page_number", strconv.Itoa(1))
		page = 1
	}
	q.Del("reset_page")
	r.URL.RawQuery = q.Encode()
	if err != nil {
		slog.Info("HandleRealizedFiltered", "err", err)
		return err
	}
	results, err := filterRealizedThesisEntries(r)
	if err != nil {
		return err
	}
	return Render(w, r, realized.SwapResults(results, page))
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
	slog.Info("HandleRealizedEntry", "entered", true)
	id_param := chi.URLParam(r, "id")
	slog.Info("HREntry", "id_param", id_param)
	thes_data, err := server.MyS.DB.RealizedThesisEntryByID(id_param)
	if err != nil {
		return err
	}
	return Render(w, r, realized.Entry(thes_data))
}

func extractRealizedThesisFromForm(r *http.Request) *types.RealizedThesisEntry {
	supHours, _ := strconv.Atoi(r.FormValue("supervisorHours"))
	assHours, _ := strconv.Atoi(r.FormValue("assistantSupervisorHours"))
	revHours, _ := strconv.Atoi(r.FormValue("reviewerHours"))
	supSettled, _ := strconv.Atoi(r.FormValue("supervisorSettled"))
	assSettled, _ := strconv.Atoi(r.FormValue("assistantSupervisorSettled"))
	revSettled, _ := strconv.Atoi(r.FormValue("reviewerSettled"))
	slog.Info("extractRealizedThesisFromForm", "supSettled", supSettled)
	slog.Info("extractRealizedThesisFromForm", "asSettled", assSettled)
	slog.Info("extractRealizedThesisFromForm", "revSettled", revSettled)
	return &types.RealizedThesisEntry{
		ThesisNumber:         strings.TrimSpace(r.FormValue("thesisNumber")),
		ExamDate:             strings.TrimSpace(r.FormValue("examDate")),
		ExamTime:             strings.TrimSpace(r.FormValue("examTime")),
		AverageStudyGrade:    strings.TrimSpace(r.FormValue("averageStudyGrade")),
		CompetencyExamGrade:  strings.TrimSpace(r.FormValue("competencyExamGrade")),
		DiplomaExamGrade:     strings.TrimSpace(r.FormValue("diplomaExamGrade")),
		FinalStudyResult:     strings.TrimSpace(r.FormValue("finalStudyResult")),
		FinalStudyResultText: strings.TrimSpace(r.FormValue("finalStudyResultText")),
		ThesisTitlePolish:    strings.TrimSpace(r.FormValue("thesisTitlePolish")),
		ThesisTitleEnglish:   strings.TrimSpace(r.FormValue("thesisTitleEnglish")),
		ThesisLanguage:       strings.TrimSpace(r.FormValue("thesisLanguage")),
		Library:              strings.TrimSpace(r.FormValue("library")),
		Student: types.Student{
			StudentNumber:  strings.TrimSpace(r.FormValue("studentNumber")),
			FirstName:      strings.TrimSpace(r.FormValue("firstNameStudent")),
			LastName:       strings.TrimSpace(r.FormValue("lastNameStudent")),
			FieldOfStudy:   strings.TrimSpace(r.FormValue("course")),
			Specialization: strings.TrimSpace(r.FormValue("specialization")),
			ModeOfStudies:  strings.TrimSpace(r.FormValue("modeOfStudies")),
			Degree:         strings.TrimSpace(r.FormValue("degree")),
		},
		ChairAcademicTitle: strings.TrimSpace(r.FormValue("chairAcademicTitle")),
		Chair: types.UniversityEmployeeEntry{
			FirstName:            strings.TrimSpace(r.FormValue("firstNameChair")),
			LastName:             strings.TrimSpace(r.FormValue("lastNameChair")),
			CurrentAcademicTitle: strings.TrimSpace(r.FormValue("chairAcademicTitle")),
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
		ReviewerAcademicTitle: strings.TrimSpace(r.FormValue("reviewerAcademicTitle")),
		Reviewer: types.UniversityEmployeeEntry{
			FirstName:            strings.TrimSpace(r.FormValue("firstNameReviewer")),
			LastName:             strings.TrimSpace(r.FormValue("lastNameReviewer")),
			CurrentAcademicTitle: strings.TrimSpace(r.FormValue("reviewerAcademicTitle")),
		},
		HourlySettlement: types.HourlySettlement{
			SupervisorHours:                 supHours,
			AssistantSupervisorHours:        assHours,
			ReviewerHours:                   revHours,
			SupervisorHoursSettled:          supSettled,
			AssistantSupervisorHoursSettled: assSettled,
			ReviewerHoursSettled:            revSettled,
		},
		Note: types.Note{
			Content: r.FormValue("thesis_note"),
		},
	}
}

func getHourlySettlementId(h types.HourlySettlement, thesis_id int) (int, error) {
	hid, err := server.MyS.DB.GetHourlySettlementIdFromRealizedThesis(thesis_id)
	if err != nil {
		return 0, err
	}
	h.Id = hid
	slog.Info("getHourlySettlementId", "hId", h.Id)
	if h.Id == 0 {
		id, err := server.MyS.DB.InsertHourlySettlement(h)
		if err != nil {
			slog.Error("getHourlySettlementId", "err", err)
			return 0, err
		}
		slog.Info("getHourlySettlementId", "inserting new HourlySettlement, id", id)
		hId := int(id)
		return hId, nil
	}
	hId, err := server.MyS.DB.UpdateHourlySettlement(h)
	if err != nil {
		slog.Error("getHourlySettlementId", "err", err)
		return 0, err
	}
	return hId, nil
}

func HandleRealizedNew(w http.ResponseWriter, r *http.Request) error {
	t := *extractRealizedThesisFromForm(r)
	t.ThesisNumber = validators.CheckThesisNumber(t.ThesisNumber, t.Student.Degree)
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
	hId, err := server.MyS.DB.InsertHourlySettlement(t.HourlySettlement)
	if err != nil {
		slog.Error("Insert HourlySettlement", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.NewEntrySwap(types.RealizedThesisEntry{}, t, errors))
	}
	t.HourlySettlement.Id = int(hId)
	tId, err := server.MyS.DB.InsertRealizedThesisByEntry(&t)
	if err != nil {
		slog.Error("InsertNew", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.NewEntrySwap(types.RealizedThesisEntry{}, t, errors))
	}
	t.Id = int(tId)
	slog.Info("thesis to db", "new_id", tId)
	slog.Info("sudent to db", "new_id", sId)
	errors.Correct = true
	return Render(w, r, realized.NewEntrySwap(t, types.RealizedThesisEntry{}, errors))
}

func HandleRealizedUpdate(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("UPDATE", "id_param", id_param)
	user, ok := r.Context().Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		slog.Error("UpdateRealized", "Failed To retrieve user", user)
	}
	t := *extractRealizedThesisFromForm(r)
	var err error
	t.Id, err = strconv.Atoi(id_param)
	slog.Info("UpdateRealizedThesis", "t before", t)
	errors, ok := validators.ValidateRealizedThesis(t)
	if user.IsAdmin {
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
		hId, err := getHourlySettlementId(t.HourlySettlement, t.Id)
		if err != nil {
			slog.Error("update HourlySettlement", "err", err)
			errors.InternalError = true
			return Render(w, r, realized.Details(t, errors))
		}
		t.HourlySettlement.Id = hId
		err = server.MyS.DB.UpdateRealizedThesisByEntry(&t)
		if err != nil {
			slog.Error("Update Thesis", "err", err)
			errors.InternalError = true
			return Render(w, r, realized.Details(t, errors))
		}
	}
	err = updateNoteRealized(t.Note.Content, t.Id, user.Id)
	if err != nil {
		slog.Error("update note", "err", err)
		errors.InternalError = true
		return Render(w, r, realized.Details(t, errors))
	}
	slog.Info("update thesis", "correct", true)
	errors.Correct = true
	slog.Info("UpdateRealizedThesis", "t after", t)
	return Render(w, r, realized.Entry(t))
}

func updateNoteRealized(newContent string, thesisId int, userId int) error {
	note, err := server.MyS.DB.GetNote(thesisId, 0, userId)
	if err != nil {
		slog.Error("updateNoteRealized", "err", err)
		return fmt.Errorf("updateNoteRealized occured error while retrieving note: %v", err)
	}
	if note.Id == 0 {
		_, err := server.MyS.DB.InsertNote(types.Note{Content: newContent, RealizedThesisID: thesisId, UniversityEmployeeID: userId})
		if err != nil {
			slog.Error("updateNoteRealized", "err", err)
			return fmt.Errorf("updateNoteRealized occured error while inserting note: %v", err)
		}
		return nil
	}
	note.Content = newContent
	err = server.MyS.DB.UpdateNote(note)
	if err != nil {
		slog.Error("updateNoteRealized", "err", err)
		return fmt.Errorf("updateNoteRealized occured error while updating note: %v", err)
	}
	return nil
}

func HandleRealizedGetNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, realized.NewEntry(types.RealizedThesisEntry{
		HourlySettlement: types.HourlySettlement{
			SupervisorHours: 10,
			ReviewerHours:   2,
		}}, types.RealizedThesisEntryErrors{}))
}

func HandleRealizedClearNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, realized.EmptySpace())
}

func HandleRealizedExcelField(w http.ResponseWriter, r *http.Request) error {
	defaultName := "Wybrane_Prace"
	return Render(w, r, components.ExcelField(defaultName))
}

func HandleRealizedClearExcelField(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, components.EmptySpace())
}
