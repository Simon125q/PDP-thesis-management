package sqlite

import (
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"thesis-management-app/types"
)

func (m *Model) AllRealizedThesis(sort_by string, desc_order bool, queryParams url.Values) ([]types.RealizedThesis, error) {
	query := fmt.Sprintf(`SELECT id, COALESCE(thesis_number, '0'), COALESCE(exam_date, '01.01.0001'), COALESCE(average_study_grade, 0), COALESCE(competency_exam_grade, 0),
    COALESCE(diploma_exam_grade, 0), COALESCE(final_study_result, ''), COALESCE(final_study_result_text, ''),
    COALESCE(thesis_title_polish, ''), COALESCE(thesis_title_english, ''), COALESCE(thesis_language, ''), COALESCE(library, ''),
    student_id, chair_id, supervisor_id, COALESCE(assistant_supervisor_id, 0), reviewer_id, COALESCE(hourly_settlement_id, 0)
    FROM Completed_Thesis`)
	query, params := m.AddSQLQueryParameters(query, queryParams)
	query = AddSQLOrder(query, sort_by, desc_order)
	slog.Info("Query", "query", query)
	slog.Info("Query", "params", params)
	rows, err := m.DB.Query(query, params...)
	if err != nil {
		return nil, err
	}
	thesis := []types.RealizedThesis{}
	for rows.Next() {
		t := types.RealizedThesis{}
		err := rows.Scan(&t.Id, &t.ThesisNumber, &t.ExamDate, &t.AverageStudyGrade, &t.CompetencyExamGrade, &t.DiplomaExamGrade,
			&t.FinalStudyResult, &t.FinalStudyResultText, &t.ThesisTitlePolish, &t.ThesisTitleEnglish, &t.ThesisLanguage,
			&t.Library, &t.StudentId, &t.ChairId, &t.SupervisorId, &t.AssistantSupervisorId, &t.ReviewerId, &t.HourlySettlementId)
		if err != nil {
			return nil, err
		}
		thesis = append(thesis, t)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return thesis, nil
}

func (m *Model) AllRealizedThesisEntries(sort_by string, desc_order bool, queryParams url.Values) ([]types.RealizedThesisEntry, error) {
	query := fmt.Sprintf(`
        SELECT 
		ct.id as thesis_id,
		COALESCE(ct.thesis_number, '0') AS thesis_number,
		COALESCE(ct.exam_date, '01.01.0001') AS exam_date,
		COALESCE(ct.average_study_grade, 0) AS average_study_grade,
		COALESCE(ct.competency_exam_grade, 0) AS competency_exam_grade,
		COALESCE(ct.diploma_exam_grade, 0) AS diploma_exam_grade,
		COALESCE(ct.final_study_result, '') AS final_study_result,
		COALESCE(ct.final_study_result_text, '') AS final_study_result_text,
		COALESCE(ct.thesis_title_polish, '') AS thesis_title_polish,
		COALESCE(ct.thesis_title_english, '') AS thesis_title_english,
		COALESCE(ct.thesis_language, '') AS thesis_language,
		COALESCE(ct.library, '') AS library,
		
		s.id AS student_id,
		COALESCE(s.student_number, '') AS student_number,
        s.first_name as student_first_name,
        s.last_name as student_last_name,
		COALESCE(s.field_of_study, '') AS student_field_of_study,
		COALESCE(s.specialization, '') AS student_specialization,
		COALESCE(s.mode_of_study, '') AS student_mode_of_study,
		
		ch.id AS chair_id,
		ch.first_name AS chair_first_name,
		ch.last_name AS chair_last_name,
		COALESCE(ch.current_academic_title, '') AS chair_curr_academic_title,
		COALESCE(ch.department_unit, '') AS chair_department_unit,
		
		sp.id AS supervisor_id,
		sp.first_name AS supervisor_first_name,
		sp.last_name AS supervisor_last_name,
		COALESCE(sp.current_academic_title, '') AS supervisor_curr_academic_title,
		COALESCE(sp.department_unit, '') AS supervisor_department_unit,
		
		COALESCE(asup.id, '0') AS assit_suppervisor_id,
		COALESCE(asup.first_name, '') AS assit_suppervisor_first_name,
		COALESCE(asup.last_name, '') AS assit_suppervisor_last_name,
		COALESCE(asup.current_academic_title, '') AS assit_suppervisor_curr_academic_title,
		COALESCE(asup.department_unit, '') AS assit_suppervisor_department_unit,
		
		r.id AS reviewer_id,
		r.first_name AS reviewer_first_name,
		r.last_name AS reviewer_last_name,
		COALESCE(r.current_academic_title, '') AS reviewer_curr_academic_title,
		COALESCE(r.department_unit, '') AS reviewer_department_unit,
		
		COALESCE(h.id, '0') AS hourly_settlement_id,
        COALESCE(h.supervisor_hours, '0'),
        COALESCE(h.assistant_supervisor_hours, '0'),
        COALESCE(h.reviewer_hours, '0')
		
	FROM 
	    Completed_Thesis ct
	LEFT JOIN Student s ON ct.student_id = s.id
	LEFT JOIN University_Employee ch ON ct.chair_id = ch.id
	LEFT JOIN University_Employee sp ON ct.supervisor_id = sp.id
	LEFT JOIN University_Employee asup ON ct.assistant_supervisor_id = asup.id
	LEFT JOIN University_Employee r ON ct.reviewer_id = r.id
	LEFT JOIN Hourly_Settlement h ON ct.hourly_settlement_id = h.id
    `)
	query, params := m.AddSQLQueryParameters(query, queryParams)
	query = AddSQLOrder(query, sort_by, desc_order)
	slog.Info("AllRealizedThesisEntries", "query", query)
	slog.Info("AllRealizedThesisEntries", "Qparams", params)
	rows, err := m.DB.Query(query, params...)
	if err != nil {
		return nil, err
	}
	thesis := []types.RealizedThesisEntry{}
	slog.Info("RealizedThesisEntry", "row", rows.Next())
	for rows.Next() {
		t := types.RealizedThesisEntry{}
		err := rows.Scan(&t.Id, &t.ThesisNumber, &t.ExamDate, &t.AverageStudyGrade, &t.CompetencyExamGrade,
			&t.DiplomaExamGrade, &t.FinalStudyResult, &t.FinalStudyResultText, &t.ThesisTitlePolish,
			&t.ThesisTitleEnglish, &t.ThesisLanguage, &t.Library,
			&t.Student.Id, &t.Student.StudentNumber, &t.Student.FirstName, &t.Student.LastName,
			&t.Student.FieldOfStudy, &t.Student.Specialization, &t.Student.ModeOfStudies,
			&t.Chair.Id, &t.Chair.FirstName, &t.Chair.LastName, &t.Chair.CurrentAcademicTitle, &t.Chair.DepartmentUnit,
			&t.Supervisor.Id, &t.Supervisor.FirstName, &t.Supervisor.LastName, &t.Supervisor.CurrentAcademicTitle, &t.Supervisor.DepartmentUnit,
			&t.AssistantSupervisor.Id, &t.AssistantSupervisor.FirstName, &t.AssistantSupervisor.LastName, &t.AssistantSupervisor.CurrentAcademicTitle, &t.AssistantSupervisor.DepartmentUnit,
			&t.Reviewer.Id, &t.Reviewer.FirstName, &t.Reviewer.LastName, &t.Reviewer.CurrentAcademicTitle, &t.Reviewer.DepartmentUnit,
			&t.HourlySettlement.Id, &t.HourlySettlement.SupervisorHours, &t.HourlySettlement.AssistantSupervisorHours, &t.HourlySettlement.ReviewerHours)
		if err != nil {
			return nil, err
		}
		thesis = append(thesis, t)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return thesis, nil
}

func (m *Model) AllRealizedThesisEntriesOld(sort_by string, desc_order bool, queryParams url.Values) ([]types.RealizedThesisEntry, error) {
	thesis, err := m.AllRealizedThesis(sort_by, desc_order, queryParams)
	if err != nil {
		return nil, err
	}
	var result []types.RealizedThesisEntry
	for _, t := range thesis {
		student, err := m.StudentById(strconv.Itoa(t.StudentId))
		if err != nil {
			return nil, err
		}
		chair, err := m.EmployeeById(strconv.Itoa(t.ChairId))
		if err != nil {
			return nil, err
		}
		supervisor, err := m.EmployeeById(strconv.Itoa(t.SupervisorId))
		if err != nil {
			return nil, err
		}
		assistant_supervisor, err := m.EmployeeById(strconv.Itoa(t.AssistantSupervisorId))
		if err != nil {
			return nil, err
		}
		reviewer, err := m.EmployeeById(strconv.Itoa(t.ReviewerId))
		if err != nil {
			return nil, err
		}
		hours, err := m.HoursById(strconv.Itoa(t.HourlySettlementId))
		if err != nil {
			return nil, err
		}
		t_entry := types.RealizedThesisEntry{
			Id:                               t.Id,
			ThesisNumber:                     t.ThesisNumber,
			ExamDate:                         t.ExamDate,
			AverageStudyGrade:                t.AverageStudyGrade,
			CompetencyExamGrade:              t.CompetencyExamGrade,
			DiplomaExamGrade:                 t.DiplomaExamGrade,
			FinalStudyResult:                 t.FinalStudyResult,
			FinalStudyResultText:             t.FinalStudyResultText,
			ThesisTitlePolish:                t.ThesisTitlePolish,
			ThesisTitleEnglish:               t.ThesisTitleEnglish,
			ThesisLanguage:                   t.ThesisLanguage,
			Library:                          t.Library,
			Student:                          student,
			ChairAcademicTitle:               t.ChairAcademicTitle,
			Chair:                            chair,
			SupervisorAcademicTitle:          t.SupervisorAcademicTitle,
			Supervisor:                       supervisor,
			AssistantSupervisorAcademicTitle: t.AssistantSupervisorAcademicTitle,
			AssistantSupervisor:              assistant_supervisor,
			ReviewerAcademicTitle:            t.ReviewerAcademicTitle,
			Reviewer:                         reviewer,
			HourlySettlement:                 hours,
		}
		result = append(result, t_entry)
	}
	return result, nil
}

func (m *Model) RealizedThesisByID(id string) (types.RealizedThesis, error) {
	query := fmt.Sprintf(`SELECT id, COALESCE(thesis_number, '0'), COALESCE(exam_date, '01.01.0001'), COALESCE(average_study_grade, 0), COALESCE(competency_exam_grade, 0),
    COALESCE(diploma_exam_grade, 0), COALESCE(final_study_result, ''), COALESCE(final_study_result_text, ''),
    COALESCE(thesis_title_polish, ''), COALESCE(thesis_title_english, ''), COALESCE(thesis_language, ''), COALESCE(library, ''),
    student_id, chair_id, supervisor_id, COALESCE(assistant_supervisor_id, 0), reviewer_id, COALESCE(hourly_settlement_id, 0)
    FROM Completed_Thesis WHERE id = %v`, id)
	slog.Info("Query by ID", "query", query)
	rows, err := m.DB.Query(query)
	if err != nil {
		return types.RealizedThesis{}, err
	}
	t := types.RealizedThesis{}
	rows.Next()
	err = rows.Scan(&t.Id, &t.ThesisNumber, &t.ExamDate, &t.AverageStudyGrade, &t.CompetencyExamGrade, &t.DiplomaExamGrade,
		&t.FinalStudyResult, &t.FinalStudyResultText, &t.ThesisTitlePolish, &t.ThesisTitleEnglish, &t.ThesisLanguage,
		&t.Library, &t.StudentId, &t.ChairId, &t.SupervisorId, &t.AssistantSupervisorId, &t.ReviewerId, &t.HourlySettlementId)
	if err != nil {
		return types.RealizedThesis{}, err
	}
	err = rows.Err()
	if err != nil {
		return types.RealizedThesis{}, err
	}
	return t, nil
}

func (m *Model) RealizedThesisEntryByID(id string) (types.RealizedThesisEntry, error) {
	t, err := m.RealizedThesisByID(id)
	student, err := m.StudentById(strconv.Itoa(t.StudentId))
	if err != nil {
		return types.RealizedThesisEntry{}, err
	}
	chair, err := m.EmployeeById(strconv.Itoa(t.ChairId))
	if err != nil {
		return types.RealizedThesisEntry{}, err
	}
	supervisor, err := m.EmployeeById(strconv.Itoa(t.SupervisorId))
	if err != nil {
		return types.RealizedThesisEntry{}, err
	}
	assistant_supervisor, err := m.EmployeeById(strconv.Itoa(t.AssistantSupervisorId))
	if err != nil {
		return types.RealizedThesisEntry{}, err
	}
	reviewer, err := m.EmployeeById(strconv.Itoa(t.ReviewerId))
	if err != nil {
		return types.RealizedThesisEntry{}, err
	}
	hours, err := m.HoursById(strconv.Itoa(t.HourlySettlementId))
	if err != nil {
		return types.RealizedThesisEntry{}, err
	}
	return types.RealizedThesisEntry{
		Id:                               t.Id,
		ThesisNumber:                     t.ThesisNumber,
		ExamDate:                         t.ExamDate,
		AverageStudyGrade:                t.AverageStudyGrade,
		CompetencyExamGrade:              t.CompetencyExamGrade,
		DiplomaExamGrade:                 t.DiplomaExamGrade,
		FinalStudyResult:                 t.FinalStudyResult,
		FinalStudyResultText:             t.FinalStudyResultText,
		ThesisTitlePolish:                t.ThesisTitlePolish,
		ThesisTitleEnglish:               t.ThesisTitleEnglish,
		ThesisLanguage:                   t.ThesisLanguage,
		Library:                          t.Library,
		Student:                          student,
		ChairAcademicTitle:               t.ChairAcademicTitle,
		Chair:                            chair,
		SupervisorAcademicTitle:          t.SupervisorAcademicTitle,
		Supervisor:                       supervisor,
		AssistantSupervisorAcademicTitle: t.AssistantSupervisorAcademicTitle,
		AssistantSupervisor:              assistant_supervisor,
		ReviewerAcademicTitle:            t.ReviewerAcademicTitle,
		Reviewer:                         reviewer,
		HourlySettlement:                 hours,
	}, nil

}
