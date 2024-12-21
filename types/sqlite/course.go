package sqlite

import (
	"fmt"
	"thesis-management-app/types"
)

func (m *Model) AllCourses() ([]types.Course, error) {
	q := `SELECT id, name FROM fields_of_study ORDER BY id DESC`
	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := []types.Course{}
	for rows.Next() {
		c := types.Course{}
		err := rows.Scan(&c.Id, &c.Name)
		if err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (m *Model) CourseById(id string) (types.Course, error) {
	if id == "0" {
		return types.Course{}, nil
	}
	query := `SELECT id, name FROM fields_of_study WHERE id = ?`
	row := m.DB.QueryRow(query, id)

	c := types.Course{}
	err := row.Scan(&c.Id, &c.Name)
	if err != nil {
		return types.Course{}, err
	}
	return c, nil
}

func (m *Model) InsertCourse(course types.Course) (int64, error) {
	query := `INSERT INTO fields_of_study (name) VALUES (?)`
	result, err := m.DB.Exec(query, course.Name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *Model) UpdateCourse(course types.Course) error {
	query := `
        UPDATE fields_of_study
        SET name = ?
        WHERE id = ?
    `
	_, err := m.DB.Exec(query, course.Name, course.Id)
	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}
	return nil
}
