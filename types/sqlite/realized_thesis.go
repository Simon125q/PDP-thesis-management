package sqlite

import (
	"fmt"
	"log/slog"
	"thesis-management-app/types"
)

func (m *Model) AllRealizedThesis(sort_by string, desc_order bool) ([]types.RealizedThesis, error) {
	order := "DESC"
	if !desc_order {
		order = "ASC"
	}
	if sort_by == "" {
		sort_by = "id"
	}
	q := fmt.Sprintf(`SELECT id, COALESCE(thesis_number, '0'), COALESCE(exam_date, '01.01.0001'), COALESCE(average_study_grade, 0), COALESCE(competency_exam_grade, 0),
    COALESCE(diploma_exam_grade, 0), COALESCE(final_study_result, ''), COALESCE(final_study_result_text, ''),
    COALESCE(thesis_title_polish, ''), COALESCE(thesis_title_english, ''), COALESCE(thesis_language, ''), COALESCE(library, ''),
    student_id, chair_id, supervisor_id, COALESCE(assistant_supervisor_id, 0), reviewer_id, COALESCE(hourly_settlement_id, 0)
    FROM Completed_Thesis ORDER BY %v %v`, sort_by, order)
	rows, err := m.DB.Query(q)
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
	slog.Info("RealizedThesis", "thesises", thesis)
	return thesis, nil
}
