package sqlite

import (
	"database/sql"
	"thesis-management-app/types"
)

type Model struct {
	DB *sql.DB
}

func (m *Model) AllUniversityEmployee() ([]types.UniversityEmployee, error) {
	q := `SELECT id, first_name, last_name, COALESCE(current_academic_title, ''), COALESCE(department_unit, '') FROM University_Employee ORDER BY id DESC`
	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	employee := []types.UniversityEmployee{}
	for rows.Next() {
		e := types.UniversityEmployee{}
		err := rows.Scan(&e.Id, &e.FirstName, &e.LastName, &e.CurrentAcademicTitle, &e.DepartmentUnit)
		if err != nil {
			return nil, err
		}
		employee = append(employee, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return employee, nil
}
