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
