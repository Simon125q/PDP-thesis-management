package sqlite

// TODO: Add all university employee current academic titles

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
	rows, err := m.DB.Query(query, params...)
	if err != nil {
		slog.Error("AllRealizedThesis", "err", err)
		return nil, err
	}
	defer rows.Close()
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
        COALESCE(chair_academic_title, '') AS chair_title,
        COALESCE(supervisor_academic_title, '') AS supervisor_title,
        COALESCE(assistant_supervisor_academic_title, '') AS assistant_title,
        COALESCE(reviewer_academic_title, '') AS reviewer_title, 
		
		s.id AS student_id,
		COALESCE(s.student_number, '') AS student_number,
        s.first_name as student_first_name,
        s.last_name as student_last_name,
		COALESCE(s.field_of_study, '') AS student_field_of_study,
		COALESCE(s.specialization, '') AS student_specialization,
		COALESCE(s.mode_of_study, '') AS student_mode_of_study,
		
		COALESCE(ch.id, '0') AS chair_id,
		COALESCE(ch.first_name, '') AS chair_first_name,
		COALESCE(ch.last_name, '') AS chair_last_name,
		COALESCE(ch.current_academic_title, '') AS chair_curr_academic_title,
		COALESCE(ch.department_unit, '') AS chair_department_unit,
		
		COALESCE(sp.id, '0') AS supervisor_id,
		COALESCE(sp.first_name, '') AS supervisor_first_name,
		COALESCE(sp.last_name, '') AS supervisor_last_name,
		COALESCE(sp.current_academic_title, '') AS supervisor_curr_academic_title,
		COALESCE(sp.department_unit, '') AS supervisor_department_unit,
		
		COALESCE(asup.id, '0') AS assit_suppervisor_id,
		COALESCE(asup.first_name, '') AS assit_suppervisor_first_name,
		COALESCE(asup.last_name, '') AS assit_suppervisor_last_name,
		COALESCE(asup.current_academic_title, '') AS assit_suppervisor_curr_academic_title,
		COALESCE(asup.department_unit, '') AS assit_suppervisor_department_unit,
		
		COALESCE(r.id, '0') AS reviewer_id,
		COALESCE(r.first_name, '') AS reviewer_first_name,
		COALESCE(r.last_name, '') AS reviewer_last_name,
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
	slog.Error("AllRealizedThesisEntries", "query", query)
	rows, err := m.DB.Query(query, params...)
	if err != nil {
		slog.Error("AllRealizedThesisEntries", "err", err)
		return nil, err
	}
	defer rows.Close()
	thesis := []types.RealizedThesisEntry{}
	for rows.Next() {
		t := types.RealizedThesisEntry{}
		err := rows.Scan(&t.Id, &t.ThesisNumber, &t.ExamDate, &t.AverageStudyGrade, &t.CompetencyExamGrade,
			&t.DiplomaExamGrade, &t.FinalStudyResult, &t.FinalStudyResultText, &t.ThesisTitlePolish,
			&t.ThesisTitleEnglish, &t.ThesisLanguage, &t.Library,
			&t.ChairAcademicTitle, &t.SupervisorAcademicTitle, &t.AssistantSupervisorAcademicTitle, &t.ReviewerAcademicTitle,
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

func (m *Model) GetAllThesisTitlesPolish(searchString string) ([]string, error) {
	query := `
        SELECT COALESCE(thesis_title_polish, '')
        FROM Completed_Thesis
        WHERE thesis_title_polish LIKE '%' || ? || '%'
    `
	rows, err := m.DB.Query(query, searchString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := []string{}

	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func (m *Model) GetAllStudentSurnames(searchString string) ([]string, error) {
	query := `
        SELECT COALESCE(last_name, '')
        FROM Student
        WHERE last_name LIKE '%' || ? || '%'
    `
	rows, err := m.DB.Query(query, searchString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := []string{}

	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func (m *Model) GetAllStudentNumbers(searchString string) ([]string, error) {
	query := `
        SELECT DISTINCT COALESCE(student_number, '')
        FROM Student
        WHERE student_number LIKE '%' || ? || '%'
    `
	rows, err := m.DB.Query(query, searchString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := []string{}

	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func (m *Model) GetAllUniversityEmployeesNames(searchString string) ([]string, error) {
    query := `
        SELECT DISTINCT COALESCE(first_name, '')
        FROM University_Employee
        WHERE first_name LIKE '%' || ? || '%'
    `
    rows, err := m.DB.Query(query, searchString)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    values := []string{}

    for rows.Next() {
        var value string
        if err := rows.Scan(&value); err != nil {
            return nil, err
        }
        values = append(values, value)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return values, nil
}

func (m *Model) GetAllUniversityEmployeesSurnames(searchString string) ([]string, error) {
    query := `
        SELECT DISTINCT COALESCE (last_name, '')
        FROM University_Employee
        WHERE last_name LIKE '%' || ? || '%'
    `
    rows, err := m.DB.Query(query, searchString)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    values := []string{}

    for rows.Next() {
        var value string
        if err := rows.Scan(&value); err != nil {
            return nil, err
        }
        values = append(values, value)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return values, nil
}



func (m *Model) GetAllUniversityEmployeesNamesAndSurnames(searchString string) ([]string, error) {
    query := `
	SELECT DISTINCT 
	COALESCE(first_name, '') || ' ' || COALESCE(last_name, '') AS result
	FROM University_Employee
	WHERE result LIKE '%' || ? || '%';`
    rows, err := m.DB.Query(query, searchString)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    values := []string{}

    for rows.Next() {
        var value string
        if err := rows.Scan(&value); err != nil {
            return nil, err
        }
        values = append(values, value)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return values, nil
}

func (m *Model) GetAllStudentsNamesAndSurnames(searchString string) ([]string, error) {
    query := `
	SELECT DISTINCT 
	COALESCE(first_name, '') || ' ' || COALESCE(last_name, '') AS result
	FROM Student
	WHERE result LIKE '%' || ? || '%';`
    rows, err := m.DB.Query(query, searchString)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    values := []string{}

    for rows.Next() {
        var value string
        if err := rows.Scan(&value); err != nil {
            return nil, err
        }
        values = append(values, value)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return values, nil
}

func (m *Model) GetAllUniversityEmployeesTitles(searchString string) ([]string, error) {
    query := `
        SELECT DISTINCT COALESCE (current_academic_title, '')
        FROM University_Employee
        WHERE current_academic_title LIKE '%' || ? || '%'
    `
    rows, err := m.DB.Query(query, searchString)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    values := []string{}

    for rows.Next() {
        var value string
        if err := rows.Scan(&value); err != nil {
            return nil, err
        }
        values = append(values, value)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return values, nil
}

func (m *Model) GetAllCourseNames(searchString string) ([]string, error) {
	query := `
        SELECT DISTINCT COALESCE(field_of_study, '')
        FROM Student
        WHERE field_of_study LIKE '%' || ? || '%'
    `
	rows, err := m.DB.Query(query, searchString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := []string{}

	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func (m *Model) RealizedThesisByID(id string) (types.RealizedThesis, error) {
	query := fmt.Sprintf(`SELECT id, COALESCE(thesis_number, '0'), COALESCE(exam_date, '01.01.0001'), COALESCE(average_study_grade, 0), COALESCE(competency_exam_grade, 0),
    COALESCE(diploma_exam_grade, 0), COALESCE(final_study_result, ''), COALESCE(final_study_result_text, ''),
    COALESCE(thesis_title_polish, ''), COALESCE(thesis_title_english, ''), COALESCE(thesis_language, ''), COALESCE(library, ''),
    COALESCE(chair_academic_title, ''), COALESCE(supervisor_academic_title, ''), COALESCE(assistant_supervisor_academic_title, ''), COALESCE(reviewer_academic_title, ''), 
    student_id, COALESCE(chair_id, '0'), COALESCE(supervisor_id, '0'), COALESCE(assistant_supervisor_id, '0'), COALESCE(reviewer_id, '0'), COALESCE(hourly_settlement_id, '0')
    FROM Completed_Thesis WHERE id = %v`, id)
	rows, err := m.DB.Query(query)
	if err != nil {
		return types.RealizedThesis{}, err
	}
	defer rows.Close()
	t := types.RealizedThesis{}
	rows.Next()
	err = rows.Scan(&t.Id, &t.ThesisNumber, &t.ExamDate, &t.AverageStudyGrade, &t.CompetencyExamGrade, &t.DiplomaExamGrade,
		&t.FinalStudyResult, &t.FinalStudyResultText, &t.ThesisTitlePolish, &t.ThesisTitleEnglish, &t.ThesisLanguage,
		&t.Library, &t.ChairAcademicTitle, &t.SupervisorAcademicTitle, &t.AssistantSupervisorAcademicTitle, &t.ReviewerAcademicTitle,
		&t.StudentId, &t.ChairId, &t.SupervisorId, &t.AssistantSupervisorId, &t.ReviewerId, &t.HourlySettlementId)
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
	if err != nil {
		return types.RealizedThesisEntry{}, err
	}
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

func (m *Model) InsertRealizedThesisByEntry(thesis *types.RealizedThesisEntry) (int64, error) {
	query := `
        INSERT INTO Completed_Thesis (
            thesis_number, exam_date, average_study_grade, competency_exam_grade,
            diploma_exam_grade, final_study_result, final_study_result_text,
            thesis_title_polish, thesis_title_english, thesis_language, library,
            chair_academic_title, supervisor_academic_title, assistant_supervisor_academic_title, reviewer_academic_title, 
            student_id, chair_id, supervisor_id, assistant_supervisor_id, reviewer_id, hourly_settlement_id
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	var sId interface{}
	if thesis.Student.Id != 0 {
		sId = thesis.Student.Id
	}
	var suId interface{}
	if thesis.Supervisor.Id != 0 {
		suId = thesis.Supervisor.Id
	}
	var asId interface{}
	if thesis.AssistantSupervisor.Id != 0 {
		asId = thesis.AssistantSupervisor.Id
	}
	var rId interface{}
	if thesis.Reviewer.Id != 0 {
		rId = thesis.Reviewer.Id
	}
	var cId interface{}
	if thesis.Chair.Id != 0 {
		cId = thesis.Chair.Id
	}
	var hId interface{}
	if thesis.HourlySettlement.Id != 0 {
		hId = thesis.HourlySettlement.Id
	}
	result, err := m.DB.Exec(query,
		thesis.ThesisNumber, thesis.ExamDate, thesis.AverageStudyGrade, thesis.CompetencyExamGrade,
		thesis.DiplomaExamGrade, thesis.FinalStudyResult, thesis.FinalStudyResultText,
		thesis.ThesisTitlePolish, thesis.ThesisTitleEnglish, thesis.ThesisLanguage, thesis.Library,
		thesis.ChairAcademicTitle, thesis.SupervisorAcademicTitle, thesis.AssistantSupervisorAcademicTitle, thesis.ReviewerAcademicTitle,
		sId, cId, suId, asId, rId, hId)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func (m *Model) UpdateRealizedThesisByEntry(thesis *types.RealizedThesisEntry) error {
	var sId interface{}
	if thesis.Student.Id != 0 {
		sId = thesis.Student.Id
	}
	var suId interface{}
	if thesis.Supervisor.Id != 0 {
		suId = thesis.Supervisor.Id
	}
	var asId interface{}
	if thesis.AssistantSupervisor.Id != 0 {
		asId = thesis.AssistantSupervisor.Id
	}
	var rId interface{}
	if thesis.Reviewer.Id != 0 {
		rId = thesis.Reviewer.Id
	}
	var cId interface{}
	if thesis.Chair.Id != 0 {
		cId = thesis.Chair.Id
	}
	var hId interface{}
	if thesis.HourlySettlement.Id != 0 {
		hId = thesis.HourlySettlement.Id
	}
	query := `UPDATE Completed_Thesis SET 
		thesis_number = ?,
		exam_date = ?,
		average_study_grade = ?,
		competency_exam_grade = ?,
		diploma_exam_grade = ?,
		final_study_result = ?,
		final_study_result_text = ?,
		thesis_title_polish = ?,
		thesis_title_english = ?,
		thesis_language = ?,
		library = ?,
        chair_academic_title = ?,
        supervisor_academic_title = ?,
        assistant_supervisor_academic_title = ?,
        reviewer_academic_title = ?, 
        student_id = ?,
        chair_id = ?,
        supervisor_id = ?,
        assistant_supervisor_id = ?,
        reviewer_id = ?,
        hourly_settlement_id = ?
	WHERE id = ?`
	_, err := m.DB.Exec(query,
		thesis.ThesisNumber,
		thesis.ExamDate,
		thesis.AverageStudyGrade,
		thesis.CompetencyExamGrade,
		thesis.DiplomaExamGrade,
		thesis.FinalStudyResult,
		thesis.FinalStudyResultText,
		thesis.ThesisTitlePolish,
		thesis.ThesisTitleEnglish,
		thesis.ThesisLanguage,
		thesis.Library,
		thesis.ChairAcademicTitle,
		thesis.SupervisorAcademicTitle,
		thesis.AssistantSupervisorAcademicTitle,
		thesis.ReviewerAcademicTitle,
		sId,
		cId,
		suId,
		asId,
		rId,
		hId,
		thesis.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to update Completed_Thesis: %w", err)
	}
	return nil
}

func (m *Model) GetStudentIdFromThesisEntry(thesisId int) (int, error) {
	query := `SELECT student_id FROM Completed_Thesis WHERE id = ?`
	rows, err := m.DB.Query(query, thesisId)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var studentId int
	rows.Next()
	err = rows.Scan(&studentId)
	if err != nil {
		return 0, err
	}
	return studentId, nil
}
