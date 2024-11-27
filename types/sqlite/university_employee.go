package sqlite

import (
	"database/sql"
	"fmt"
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

func (m *Model) EmployeeById(id string) (types.UniversityEmployee, error) {
	if id == "0" {
		return types.UniversityEmployee{}, nil
	}
	query := fmt.Sprintf(`SELECT id, first_name, last_name,
    COALESCE(current_academic_title, ''), COALESCE(department_unit, '')
    FROM University_Employee WHERE id = %v`, id)
	rows, err := m.DB.Query(query)
	if err != nil {
		return types.UniversityEmployee{}, err
	}
	e := types.UniversityEmployee{}
	rows.Next()
	err = rows.Scan(&e.Id, &e.FirstName, &e.LastName, &e.CurrentAcademicTitle, &e.DepartmentUnit)
	if err != nil {
		return types.UniversityEmployee{}, err
	}
	err = rows.Err()
	if err != nil {
		return types.UniversityEmployee{}, err
	}
	return e, nil
}
