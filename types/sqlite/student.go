package sqlite

// TODO: add degree

import (
	"fmt"
	"thesis-management-app/types"
)

func (m *Model) StudentById(id string) (types.Student, error) {
	if id == "0" {
		return types.Student{}, nil
	}
	query := fmt.Sprintf(`SELECT id, COALESCE(student_number, '0'), first_name, last_name,
    COALESCE(field_of_study, ''), COALESCE(specialization, ''), COALESCE(mode_of_study, '')
    FROM Student WHERE id = %v`, id)
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

func (m *Model) StudentByNumber(studentNumber string) (types.Student, error) {
	query := fmt.Sprintf(`SELECT id, COALESCE(student_number, '0'), first_name, last_name,
    COALESCE(field_of_study, ''), COALESCE(specialization, ''), COALESCE(mode_of_study, ''))
    FROM Student WHERE student_number = %v`, studentNumber)
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

func (m *Model) InsertStudent(student types.Student) (int64, error) {
	query := `
        INSERT INTO Student (student_number, first_name, last_name, field_of_study, specialization, mode_of_study)
        VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := m.DB.Exec(query, student.StudentNumber, student.FirstName, student.LastName, student.FieldOfStudy, student.Specialization, student.ModeOfStudies)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
