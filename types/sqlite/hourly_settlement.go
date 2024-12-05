package sqlite

import (
	"fmt"
	"thesis-management-app/types"
)

func (m *Model) HoursById(id string) (types.HourlySettlement, error) {
	if id == "0" {
		return types.HourlySettlement{}, nil
	}
	query := fmt.Sprintf(`SELECT id, supervisor_hours, 
        assistant_supervisor_hours, reviewer_hours)
    FROM Hourly_Settlement WHERE id = %v`, id)
	rows, err := m.DB.Query(query)
	if err != nil {
		return types.HourlySettlement{}, err
	}
	defer rows.Close()
	h := types.HourlySettlement{}
	rows.Next()
	err = rows.Scan(&h.Id, &h.SupervisorHours, &h.AssistantSupervisorHours, &h.ReviewerHours)
	if err != nil {
		return types.HourlySettlement{}, err
	}
	err = rows.Err()
	if err != nil {
		return types.HourlySettlement{}, err
	}
	return h, nil
}
