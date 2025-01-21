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
	q := `SELECT id, first_name, last_name, COALESCE(current_academic_title, ''), COALESCE(department_unit, '') FROM University_Employee ORDER BY UPPER(last_name) ASC`
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

func (m *Model) GetSortedEmps(sortBy string, order string, searchTerm string) ([]types.UniversityEmployeeEntry, error) {
	if sortBy == "" {
		sortBy = "name"
	}
	if order == "" {
		order = "DESC"
	}
	if searchTerm == "" {
		searchTerm = "%"
	}

	q := fmt.Sprintf("SELECT id, first_name, last_name, COALESCE(current_academic_title, '') as current_academic_title, COALESCE(department_unit, '') AS department_unit FROM University_Employee WHERE last_name LIKE '%%%s%%' ORDER BY UPPER(%s) %s", searchTerm, sortBy, order)

	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	emps := []types.UniversityEmployeeEntry{}
	for rows.Next() {
		emp := types.UniversityEmployeeEntry{}
		err := rows.Scan(&emp.Id, &emp.FirstName, &emp.LastName, &emp.CurrentAcademicTitle, &emp.DepartmentUnit)
		if err != nil {
			return nil, err
		}
		emps = append(emps, emp)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return emps, nil
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
	var totalThesisCount int

	err := m.DB.QueryRow(`
		SELECT 
			(SELECT COUNT(*) FROM Completed_Thesis
			WHERE supervisor_id = ? OR assistant_supervisor_id = ? OR reviewer_id = ?) +
			(SELECT COUNT(*) FROM Thesis_To_Be_Completed
			WHERE supervisor_id = ? OR assistant_supervisor_id = ?)`, id, id, id, id, id).
		Scan(&totalThesisCount)

	if err != nil {
		return "error", err
	}

	return fmt.Sprintf("%d", totalThesisCount), nil
}

func (m *Model) EmployeeIdByName(name string) (int, error) {
	query := `SELECT id FROM University_Employee WHERE first_name || ' ' || last_name = ?`
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

	var totalThesisCount int
	err = m.DB.QueryRow(`
		SELECT 
			(SELECT COUNT(*) FROM Completed_Thesis WHERE supervisor_id = ? OR assistant_supervisor_id = ? OR reviewer_id = ?) +
			(SELECT COUNT(*) FROM Thesis_To_Be_Completed WHERE supervisor_id = ? OR assistant_supervisor_id = ?)`,
		id, id, id, id, id).Scan(&totalThesisCount)
	if err != nil {
		return types.UniversityEmployeeEntry{}, err
	}

	return types.UniversityEmployeeEntry{
		Id:                   empl.Id,
		FirstName:            empl.FirstName,
		LastName:             empl.LastName,
		CurrentAcademicTitle: empl.CurrentAcademicTitle,
		DepartmentUnit:       empl.DepartmentUnit,
		ThesisCount:          fmt.Sprintf("%d", totalThesisCount),
	}, nil
}

func (m *Model) GetPersonIDByFullName(name string) ([]string, error) {
	queryString := "SELECT id FROM University_Employee WHERE CONCAT(first_name, ' ', last_name) LIKE ?"
	slog.Info(queryString)
	rows, err := m.DB.Query(queryString, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var str string
		rows.Scan(&str)
		ids = append(ids, str)
	}
	return ids, nil
}
