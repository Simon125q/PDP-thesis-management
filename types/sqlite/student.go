package sqlite

import (
	"fmt"
	"log/slog"
	"thesis-management-app/types"
)

func (m *Model) StudentById(id string) (types.Student, error) {
	if id != "0" {
		return types.Student{}, nil
	}
	query := fmt.Sprintf(`SELECT id, COALESCE(student_number, '0'), first_name, last_name,
    COALESCE(field_of_study, ''), COALESCE(specialization, ''), COALESCE(mode_of_studies, '')
    FROM Student WHERE id = %v`, id)
	slog.Info("Student by ID", "query", query)
	rows, err := m.DB.Query(query)
	if err != nil {
		return types.Student{}, err
	}
	s := types.Student{}
	rows.Next()
	err = rows.Scan(&s.Id, &s.StudentNumber, &s.FirstName, &s.LastName, &s.FieldOfStudy,
		&s.Specialization, &s.ModeOfStudies)
	if err != nil {
		return types.Student{}, err
	}
	err = rows.Err()
	if err != nil {
		return types.Student{}, err
	}
	return s, nil
}
