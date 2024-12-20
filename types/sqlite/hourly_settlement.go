package sqlite

import (
	"fmt"
	"log/slog"
	"thesis-management-app/types"
)

func (m *Model) HoursById(id string) (types.HourlySettlement, error) {
	if id == "0" {
		return types.HourlySettlement{SupervisorHours: 10, ReviewerHours: 2}, nil
	}
	query := fmt.Sprintf(`SELECT id,
        COALESCE(supervisor_hours, '10'), COALESCE(assistant_supervisor_hours, '0'), COALESCE(reviewer_hours, '2'),
        COALESCE(is_supervisor_settled, '0'), COALESCE(is_assistant_supervisor_settled, '0'), COALESCE(is_reviewer_settled, '0')
    FROM Hourly_Settlement WHERE id = ?`)
	rows, err := m.DB.Query(query, id)
	if err != nil {
		return types.HourlySettlement{}, err
	}
	defer rows.Close()
	h := types.HourlySettlement{}
	rows.Next()
	err = rows.Scan(&h.Id, &h.SupervisorHours, &h.AssistantSupervisorHours,
		&h.ReviewerHours, &h.SupervisorHoursSettled, &h.AssistantSupervisorHoursSettled, &h.ReviewerHoursSettled)
	if err != nil {
		return types.HourlySettlement{}, err
	}
	err = rows.Err()
	if err != nil {
		return types.HourlySettlement{}, err
	}
	return h, nil
}

func (m *Model) InsertHourlySettlement(hours types.HourlySettlement) (int64, error) {
	slog.Info("InsertHourlySettlement", "hours", hours)
	query := `
        INSERT INTO Hourly_Settlement (
        supervisor_hours, assistant_supervisor_hours, reviewer_hours,
        is_supervisor_settled, is_assistant_supervisor_settled, is_reviewer_settled
    )
        VALUES (?, ?, ?, ?, ?, ?)`
	result, err := m.DB.Exec(query, hours.SupervisorHours, hours.AssistantSupervisorHours, hours.ReviewerHours,
		hours.SupervisorHoursSettled, hours.AssistantSupervisorHoursSettled, hours.ReviewerHoursSettled)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *Model) UpdateHourlySettlement(hours types.HourlySettlement) (int, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Ensure rollback if commit doesn't happen
	query := `
    UPDATE Hourly_Settlement SET supervisor_hours = ?, assistant_supervisor_hours = ?, reviewer_hours= ?,
        is_supervisor_settled = ?, is_assistant_supervisor_settled = ?, is_reviewer_settled = ?
    WHERE id = ?`
	_, err = tx.Exec(query, hours.SupervisorHours, hours.AssistantSupervisorHours, hours.ReviewerHours,
		hours.SupervisorHoursSettled, hours.AssistantSupervisorHoursSettled, hours.ReviewerHoursSettled, hours.Id)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return hours.Id, nil
}
