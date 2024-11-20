package sqlite

import (
	"fmt"
	"log/slog"
	"net/url"
	"thesis-management-app/types"
)

func (m *Model) AllRealizedThesis(sort_by string, desc_order bool, queryParams url.Values) ([]types.RealizedThesis, error) {
	query := fmt.Sprintf(`SELECT id, COALESCE(thesis_number, '0'), COALESCE(exam_date, '01.01.0001'), COALESCE(average_study_grade, 0), COALESCE(competency_exam_grade, 0),
    COALESCE(diploma_exam_grade, 0), COALESCE(final_study_result, ''), COALESCE(final_study_result_text, ''),
    COALESCE(thesis_title_polish, ''), COALESCE(thesis_title_english, ''), COALESCE(thesis_language, ''), COALESCE(library, ''),
    student_id, chair_id, supervisor_id, COALESCE(assistant_supervisor_id, 0), reviewer_id, COALESCE(hourly_settlement_id, 0)
    FROM Completed_Thesis`)
	query, params := AddSQLQueryParameters(query, queryParams)
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
