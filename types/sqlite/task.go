package sqlite

import (
	"fmt"
	"log/slog"
	"thesis-management-app/types"
)

func (m *Model) GetTasksByThesisId(thesis_id int) ([]types.Task, error) {
	query := `SELECT id, content, COALESCE(is_completed, 0) FROM Task WHERE thesis_to_be_completed_id = ?`
	rows, err := m.DB.Query(query, thesis_id)
	if err != nil {
		return nil, fmt.Errorf("GetTasksByThesisId -> %v", err)
	}
	defer rows.Close()
	tasks := []types.Task{}
	for rows.Next() {
		t := types.Task{}
		err := rows.Scan(&t.Id, &t.Content, &t.IsCompleted)
		if err != nil {
			return nil, fmt.Errorf("GetTasksByThesisId -> %v", err)
		}
		tasks = append(tasks, t)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetTasksByThesisId -> %v", err)
	}
	return tasks, nil
}

func (m *Model) InsertTask(task types.Task) (int64, error) {
	slog.Info("InsertTask", "task", task)
	query := `
        INSERT INTO Task(
        content, thesis_to_be_completed_id
    )
        VALUES (?, ?)`
	result, err := m.DB.Exec(query, task.Content, task.OngoingThesisID)
	if err != nil {
		return 0, fmt.Errorf("InsertTask -> %v", err)
	}
	return result.LastInsertId()
}

func (m *Model) UpdateTaskCompletnes(taskId, isChecked int) error {
	query := `
        UPDATE Task
        SET 
            is_completed = ?
        WHERE id = ?
    `
	slog.Info("UpdateTask", "TaskId", taskId)

	_, err := m.DB.Exec(query,
		isChecked,
		taskId,
	)
	if err != nil {
		return fmt.Errorf("UpdateTask -> %w", err)
	}
	return nil
}

func (m *Model) UpdateTask(task types.Task) error {
	query := `
        UPDATE Task
        SET 
            content = ?,
            is_completed = ?
        WHERE id = ?
    `
	slog.Info("UpdateTask", "Task", task)

	_, err := m.DB.Exec(query,
		task.Content,
		task.IsCompleted,
		task.Id,
	)
	if err != nil {
		return fmt.Errorf("UpdateTask -> %w", err)
	}
	return nil
}

func (m *Model) CheckIfAllTaskAreCompleted(thesisId int) bool {
	//TODO: Check if all tasks conected with thesisId are completed
	return true
}
