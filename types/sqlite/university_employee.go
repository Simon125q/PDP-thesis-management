package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
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
	defer rows.Close()
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

func (m *Model) AllUniversityEmployeeEntries(sort_by string, desc_order bool, queryParams url.Values) ([]types.UniversityEmployeeEntry, error) {
	query := fmt.Sprintf(`
        SELECT 
            e.id AS employee_id,
            COALESCE(e.first_name, '') AS first_name,
            COALESCE(e.last_name, '') AS last_name,
            COALESCE(e.current_academic_title, '') AS current_academic_title,
            COALESCE(e.department_unit, '') AS department_unit
        FROM 
            University_Employee e
    `)

	query, params := m.AddSQLQueryParameters(query, queryParams)
	query = AddSQLOrder(query, sort_by, desc_order)

	slog.Info("AllUniversityEmployeeEntries", "query", query)
	slog.Info("AllUniversityEmployeeEntries", "params", params)

	rows, err := m.DB.Query(query, params...)
	if err != nil {
		slog.Error("AllUniversityEmployeeEntries", "err", err)
		return nil, err
	}
	defer rows.Close()

	employees := []types.UniversityEmployeeEntry{}
	for rows.Next() {
		e := types.UniversityEmployeeEntry{}
		err := rows.Scan(&e.Id, &e.FirstName, &e.LastName, &e.CurrentAcademicTitle, &e.DepartmentUnit)
		if err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (m *Model) EmployeeById(id string) (types.UniversityEmployeeEntry, error) {
	if id == "0" {
		return types.UniversityEmployeeEntry{}, nil
	}
	query := fmt.Sprintf(`SELECT id, first_name, last_name,
    COALESCE(current_academic_title, ''), COALESCE(department_unit, '')
    FROM University_Employee WHERE id = ?`)
	rows, err := m.DB.Query(query, id)
	if err != nil {
		return types.UniversityEmployeeEntry{}, err
	}
	defer rows.Close()
	e := types.UniversityEmployeeEntry{}
	rows.Next()
	err = rows.Scan(&e.Id, &e.FirstName, &e.LastName, &e.CurrentAcademicTitle, &e.DepartmentUnit)
	if err != nil {
		return types.UniversityEmployeeEntry{}, err
	}
	err = rows.Err()
	if err != nil {
		return types.UniversityEmployeeEntry{}, err
	}
	return e, nil
}

func (m *Model) ThesisCountByEmpId(id string) (string, error) {
	var thesisCount string
	err := m.DB.QueryRow(`
	    SELECT COUNT(*)
	    FROM Completed_Thesis
	    WHERE supervisor_id = ? OR assistant_supervisor_id = ? OR reviewer_id = ?`, id, id, id).Scan(&thesisCount)
	if err != nil {
		return "error", err
	}
	return thesisCount, nil
}

func (m *Model) EmployeeIdByName(name string) (int, error) {
	query := fmt.Sprintf(`SELECT id FROM University_Employee WHERE first_name || ' ' || last_name = ?`)
	slog.Info("employeeId by Name", "q", query)
	slog.Info("employeeId by Name", "name", name)
	row := m.DB.QueryRow(query, name)
	id := 0
	err := row.Scan(&id)
	slog.Info("employeeId by Name", "id", id)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *Model) InsertUniversityEmployee(employee types.UniversityEmployeeEntry) (int64, error) {
	query := `
        INSERT INTO University_Employee (first_name, last_name, current_academic_title, department_unit)
        VALUES (?, ?, ?, ?)`
	result, err := m.DB.Exec(query, employee.FirstName, employee.LastName, employee.CurrentAcademicTitle, employee.DepartmentUnit)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *Model) UpdateEmployee(empl *types.UniversityEmployeeEntry) error {
	query := `
        UPDATE University_Employee
        SET 
            first_name = ?,
            last_name = ?,
            current_academic_title = ?,
            department_unit = ?
        WHERE id = ?
    `
	slog.Info("Executing UpdateEmployee Query",
		"query", query,
		"params", empl,
	)

	_, err := m.DB.Exec(query,
		empl.FirstName,
		empl.LastName,
		empl.CurrentAcademicTitle,
		empl.DepartmentUnit,
		empl.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to update university_employee: %w", err)
	}
	return nil
}

func (m *Model) EmployeeEntryByID(id string) (types.UniversityEmployeeEntry, error) {
	empl, err := m.EmployeeById(id)
	if err != nil {
		return types.UniversityEmployeeEntry{}, err
	}

	var thesisCount int
	err = m.DB.QueryRow(`
    SELECT COUNT(*)
    FROM Completed_Thesis
    WHERE supervisor_id = ? OR assistant_supervisor_id = ? OR reviewer_id = ?`, id, id, id).Scan(&thesisCount)
	if err != nil {
		return types.UniversityEmployeeEntry{}, err
	}

	thesisCount = 1

	return types.UniversityEmployeeEntry{
		Id:                   empl.Id,
		FirstName:            empl.FirstName,
		LastName:             empl.LastName,
		CurrentAcademicTitle: empl.CurrentAcademicTitle,
		DepartmentUnit:       empl.DepartmentUnit,
		ThesisCount:          fmt.Sprintf("%d", thesisCount),
	}, nil
}
