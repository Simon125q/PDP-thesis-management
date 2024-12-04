package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
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

func (m *Model) EmployeeIdByName(name string) (int, error) {

	query := fmt.Sprintf(`SELECT id FROM University_Employee WHERE first_name || ' ' || last_name = ?`)
	slog.Info("employeeId by Name", "q", query)
	row := m.DB.QueryRow(query, name)
	id := 0
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *Model) InsertUniversityEmployee(employee types.UniversityEmployee) (int64, error) {
	query := `
        INSERT INTO University_Employee (first_name, last_name, current_academic_title, department_unit)
        VALUES (?, ?, ?, ?)`
	result, err := m.DB.Exec(query, employee.FirstName, employee.LastName, employee.CurrentAcademicTitle, employee.DepartmentUnit)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
