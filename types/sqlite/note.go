package sqlite

import (
	"fmt"
	"log/slog"
	"thesis-management-app/types"
)

type Note struct {
	Id                   int
	Content              string
	UniversityEmployeeID int
	RealizedThesisID     int
	OngoingThesisID      int
}

func (m *Model) GetNote(realizedThesId, ongoingThesisId, userId int) (types.Note, error) {
	slog.Info("GetNote", "userId", userId)
	slog.Info("GetNote", "realizedThesId", realizedThesId)
	slog.Info("GetNote", "ongoingThesisId", ongoingThesisId)
	var query string
	var thesis_id int
	if realizedThesId == 0 {
		query = `SELECT id, content FROM Note WHERE thesis_to_be_completed_id = ? AND university_employee_id = ?`
		thesis_id = ongoingThesisId
	} else {
		query = `SELECT id, content FROM Note WHERE completed_thesis_id = ? AND university_employee_id = ?`
		thesis_id = realizedThesId
	}
	rows, err := m.DB.Query(query, thesis_id, userId)
	if err != nil {
		return types.Note{}, err
	}
	defer rows.Close()
	resultNote := types.Note{}
	if !rows.Next() {
		return types.Note{}, nil
	}
	err = rows.Scan(&resultNote.Id, &resultNote.Content)
	if err != nil {
		return types.Note{}, err
	}
	return resultNote, nil
}

func (m *Model) InsertNote(note types.Note) (int64, error) {
	slog.Info("InsertNote", "note", note)
	if note.RealizedThesisID == 0 {
		query := `
            INSERT INTO Note(
            content, thesis_to_be_completed_id, university_employee_id
        )
            VALUES (?, ?, ?)`
		result, err := m.DB.Exec(query, note.Content, note.OngoingThesisID, note.UniversityEmployeeID)
		if err != nil {
			return 0, err
		}
		return result.LastInsertId()
	}
	query := `
        INSERT INTO Note(
        content, completed_thesis_id, university_employee_id
    )
        VALUES (?, ?, ?)`
	result, err := m.DB.Exec(query, note.Content, note.RealizedThesisID, note.UniversityEmployeeID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *Model) UpdateNote(note types.Note) error {
	query := `
        UPDATE Note
        SET 
            content = ?
        WHERE id = ?
    `
	slog.Info("UpdateNote", "note", note)

	_, err := m.DB.Exec(query,
		note.Content,
		note.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}
	return nil
}
