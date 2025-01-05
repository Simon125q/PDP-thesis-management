package sqlite

import (
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"thesis-management-app/types"
)

func (m *Model) AllOngoingThesisEntries(sort_by string, desc_order bool, page_num, page_limit int, queryParams url.Values) ([]types.OngoingThesisEntry, error) {
	query := `
        SELECT 
		ct.id as thesis_id,
		COALESCE(ct.thesis_number, '0') AS thesis_number,
		COALESCE(ct.topic_polish, '') AS thesis_title_polish,
		COALESCE(ct.topic_english, '') AS thesis_title_english,
		COALESCE(ct.thesis_language, '') AS thesis_language,
        COALESCE(ct.supervisor_academic_title, '') AS supervisor_title,
        COALESCE(ct.assistant_supervisor_academic_title, '') AS assistant_title,
        COALESCE(ct.topic_scan, 'false') AS archived,
		
		s.id AS student_id,
		COALESCE(s.student_number, '') AS student_number,
        s.first_name as student_first_name,
        s.last_name as student_last_name,
		COALESCE(s.field_of_study, '') AS student_field_of_study,
        COALESCE(s.degree, '') AS student_degree,
		COALESCE(s.specialization, '') AS student_specialization,
		COALESCE(s.mode_of_study, '') AS student_mode_of_study,
		
		COALESCE(sp.id, '0') AS supervisor_id,
		COALESCE(sp.first_name, '') AS supervisor_first_name,
		COALESCE(sp.last_name, '') AS supervisor_last_name,
		COALESCE(sp.current_academic_title, '') AS supervisor_curr_academic_title,
		COALESCE(sp.department_unit, '') AS supervisor_department_unit,
		
		COALESCE(asup.id, '0') AS assit_suppervisor_id,
		COALESCE(asup.first_name, '') AS assit_suppervisor_first_name,
		COALESCE(asup.last_name, '') AS assit_suppervisor_last_name,
		COALESCE(asup.current_academic_title, '') AS assit_suppervisor_curr_academic_title,
		COALESCE(asup.department_unit, '') AS assit_suppervisor_department_unit

	FROM 
	    Thesis_To_Be_Completed ct
	LEFT JOIN Student s ON ct.student_id = s.id
	LEFT JOIN University_Employee sp ON ct.supervisor_id = sp.id
	LEFT JOIN University_Employee asup ON ct.assistant_supervisor_id = asup.id
    `
	query, params := m.AddSQLQueryParameters(query, queryParams)
	query = AddSQLOrder(query, sort_by, desc_order)
	query = AddSQLPagination(query, page_num, page_limit)
	slog.Info("AllOngoingThesisEntries", "query", query)
	slog.Info("AllOngoingThesisEntries", "params", params)
	rows, err := m.DB.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("AllOngoingThesisEntries query error %v", err)
	}
	defer rows.Close()
	thesis := []types.OngoingThesisEntry{}
	for rows.Next() {
		t := types.OngoingThesisEntry{}
		err := rows.Scan(&t.Id, &t.ThesisNumber, &t.ThesisTitlePolish, &t.ThesisTitleEnglish,
			&t.ThesisLanguage, &t.SupervisorAcademicTitle, &t.AssistantSupervisorAcademicTitle, &t.Archived,
			&t.Student.Id, &t.Student.StudentNumber, &t.Student.FirstName, &t.Student.LastName,
			&t.Student.FieldOfStudy, &t.Student.Degree, &t.Student.Specialization, &t.Student.ModeOfStudies,
			&t.Supervisor.Id, &t.Supervisor.FirstName, &t.Supervisor.LastName, &t.Supervisor.CurrentAcademicTitle, &t.Supervisor.DepartmentUnit,
			&t.AssistantSupervisor.Id, &t.AssistantSupervisor.FirstName, &t.AssistantSupervisor.LastName,
			&t.AssistantSupervisor.CurrentAcademicTitle, &t.AssistantSupervisor.DepartmentUnit)
		if err != nil {
			return nil, err
		}
		thesis = append(thesis, t)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("AllOngoingThesisEntries rows error %v", err)
	}
	return thesis, nil
}

func (m *Model) OngoingThesisEntryByID(id string) (types.OngoingThesisEntry, error) {
	t, err := m.OngoingThesisByID(id)
	if err != nil {
		return types.OngoingThesisEntry{}, fmt.Errorf("OngoingThesisEntryByID -> %v", err)
	}
	student, err := m.StudentById(strconv.Itoa(t.StudentId))
	if err != nil {
		return types.OngoingThesisEntry{}, fmt.Errorf("OngoingThesisEntryByID -> %v", err)
	}
	supervisor, err := m.EmployeeById(strconv.Itoa(t.SupervisorId))
	if err != nil {
		return types.OngoingThesisEntry{}, fmt.Errorf("OngoingThesisEntryByID -> %v", err)
	}
	assistant_supervisor, err := m.EmployeeById(strconv.Itoa(t.AssistantSupervisorId))
	if err != nil {
		return types.OngoingThesisEntry{}, fmt.Errorf("OngoingThesisEntryByID -> %v", err)
	}
	return types.OngoingThesisEntry{
		Id:                               t.Id,
		ThesisNumber:                     t.ThesisNumber,
		ThesisTitlePolish:                t.ThesisTitlePolish,
		ThesisTitleEnglish:               t.ThesisTitleEnglish,
		ThesisLanguage:                   t.ThesisLanguage,
		Student:                          student,
		SupervisorAcademicTitle:          t.SupervisorAcademicTitle,
		Supervisor:                       supervisor,
		AssistantSupervisorAcademicTitle: t.AssistantSupervisorAcademicTitle,
		AssistantSupervisor:              assistant_supervisor,
		Archived:                         t.Archived,
	}, nil

}

func (m *Model) OngoingThesisByID(id string) (types.OngoingThesis, error) {
	query := `SELECT id, COALESCE(thesis_number, '0'), 
    COALESCE(topic_polish, ''), COALESCE(topic_english, ''), COALESCE(thesis_language, ''), 
    COALESCE(supervisor_academic_title, ''), COALESCE(assistant_supervisor_academic_title, ''), 
    student_id, COALESCE(supervisor_id, '0'), COALESCE(assistant_supervisor_id, '0'), COALESCE(topic_scan, 'false') 
    FROM Thesis_To_Be_Completed WHERE id = ?`
	rows, err := m.DB.Query(query, id)
	if err != nil {
		return types.OngoingThesis{}, fmt.Errorf("OngoingThesisByID query error -> %v", err)
	}
	defer rows.Close()
	t := types.OngoingThesis{}
	rows.Next()
	err = rows.Scan(&t.Id, &t.ThesisNumber,
		&t.ThesisTitlePolish, &t.ThesisTitleEnglish, &t.ThesisLanguage,
		&t.SupervisorAcademicTitle, &t.AssistantSupervisorAcademicTitle,
		&t.StudentId, &t.SupervisorId, &t.AssistantSupervisorId, &t.Archived)
	if err != nil {
		return types.OngoingThesis{}, fmt.Errorf("OngoingThesisByID scan error -> %v", err)
	}
	err = rows.Err()
	if err != nil {
		return types.OngoingThesis{}, fmt.Errorf("OngoingThesisByID rows error -> %v", err)
	}
	slog.Info("OngoingThesisByID", "thesis", t)
	return t, nil
}

func (m *Model) InsertOngoingThesisByEntry(thesis *types.OngoingThesisEntry) (int64, error) {
	query := `
        INSERT INTO Thesis_To_Be_Completed (
            thesis_number,
            topic_polish, topic_english, thesis_language, 
            supervisor_academic_title, assistant_supervisor_academic_title, 
            student_id, supervisor_id, assistant_supervisor_id
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
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
	result, err := m.DB.Exec(query,
		thesis.ThesisNumber,
		thesis.ThesisTitlePolish, thesis.ThesisTitleEnglish, thesis.ThesisLanguage,
		thesis.SupervisorAcademicTitle, thesis.AssistantSupervisorAcademicTitle,
		sId, suId, asId)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *Model) UpdateOngoingThesisByEntry(thesis *types.OngoingThesisEntry) error {
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
	query := `UPDATE Thesis_To_Be_Completed SET 
		thesis_number = ?,
		topic_polish = ?,
		topic_english = ?,
		thesis_language = ?,
        supervisor_academic_title = ?,
        assistant_supervisor_academic_title = ?,
        student_id = ?,
        supervisor_id = ?,
        assistant_supervisor_id = ?
	WHERE id = ?`
	_, err := m.DB.Exec(query,
		thesis.ThesisNumber,
		thesis.ThesisTitlePolish,
		thesis.ThesisTitleEnglish,
		thesis.ThesisLanguage,
		thesis.SupervisorAcademicTitle,
		thesis.AssistantSupervisorAcademicTitle,
		sId,
		suId,
		asId,
		thesis.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to update Thesis_To_Be_Completed: %w", err)
	}
	return nil
}

func (m *Model) ArchiveOngoingThesis(thesisId int) error {
	query := `UPDATE Thesis_To_Be_Completed SET 
		topic_scan = ?
	WHERE id = ?`
	_, err := m.DB.Exec(query,
		"true",
		thesisId,
	)
	if err != nil {
		return fmt.Errorf("failed to archive Thesis_To_Be_Completed: %w", err)
	}
	return nil
}
