package sqlite

import (
	"fmt"
	"log/slog"
	"net/url"
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
		err := rows.Scan(&t.Id, &t.ThesisNumber, &t.ThesisTitlePolish,
			&t.ThesisTitleEnglish, &t.ThesisLanguage, &t.SupervisorAcademicTitle, &t.AssistantSupervisorAcademicTitle,
			&t.Student.Id, &t.Student.StudentNumber, &t.Student.FirstName, &t.Student.LastName,
			&t.Student.FieldOfStudy, &t.Student.Degree, &t.Student.Specialization, &t.Student.ModeOfStudies,
			&t.Supervisor.Id, &t.Supervisor.FirstName, &t.Supervisor.LastName, &t.Supervisor.CurrentAcademicTitle, &t.Supervisor.DepartmentUnit,
			&t.AssistantSupervisor.Id, &t.AssistantSupervisor.FirstName, &t.AssistantSupervisor.LastName, &t.AssistantSupervisor.CurrentAcademicTitle, &t.AssistantSupervisor.DepartmentUnit)
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
