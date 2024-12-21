package sqlite

import (
	"fmt"
	"thesis-management-app/types"
)

func (m *Model) AllSpecializations() ([]types.Specialization, error) {
	q := `SELECT id, name FROM Specializations ORDER BY id DESC`
	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	specs := []types.Specialization{}
	for rows.Next() {
		s := types.Specialization{}
		err := rows.Scan(&s.Id, &s.Name)
		if err != nil {
			return nil, err
		}
		specs = append(specs, s)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return specs, nil
}

func (m *Model) SpecializationById(id string) (types.Specialization, error) {
	if id == "0" {
		return types.Specialization{}, nil
	}
	query := `SELECT id, name FROM Specializations WHERE id = ?`
	row := m.DB.QueryRow(query, id)

	s := types.Specialization{}
	err := row.Scan(&s.Id, &s.Name)
	if err != nil {
		return types.Specialization{}, err
	}
	return s, nil
}

func (m *Model) InsertSpecialization(spec types.Specialization) (int64, error) {
	query := `INSERT INTO Specializations (name) VALUES (?)`
	result, err := m.DB.Exec(query, spec.Name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *Model) UpdateSpecialization(spec types.Specialization) error {
	query := `
        UPDATE Specializations
        SET name = ?
        WHERE id = ?
    `
	_, err := m.DB.Exec(query, spec.Name, spec.Id)
	if err != nil {
		return fmt.Errorf("failed to update specialization: %w", err)
	}
	return nil
}
