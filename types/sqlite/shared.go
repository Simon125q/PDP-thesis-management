package sqlite

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

func (m *Model) GetSupervisorID(name string) ([]string, error) {
	queryString := "SELECT id FROM 'University_Employee' WHERE "
	strSplit := strings.Split(name, " ")
	var parts []string
	parts = append(parts, strSplit[0])
	if len(strSplit) > 1 {
		combined := ""
		for str := range strSplit {
			if str == 0 {
				continue
			}
			combined = combined + strSplit[str] + " "
		}
		if combined != "" {
			combined = combined[:len(combined)-1]
		}
		parts = append(parts, combined)
	}
	switch len(parts) {
	case 1:
		queryString = queryString + "first_name = '" + parts[0] + "' OR last_name = '" + parts[0] + "'"
		break
	case 2:
		queryString = queryString + "first_name = '" + parts[0] + "' AND last_name = '" + parts[1] + "'"
		break
	}
	rows, err := m.DB.Query(queryString)
	if err != nil {
		return nil, err
	}
	var ids []string
	for rows.Next() {
		var str string
		rows.Scan(&str)
		ids = append(ids, str)
	}
	return ids, nil
}

func (m *Model) AddSQLQueryParameters(baseQuery string, params url.Values) (string, []interface{}) {
	var conditions []string
	var values []interface{}

	for key, value := range params {
		//str := "Key: " + key + " Val: " + value[0]
		//slog.Info(str)
		switch key {
		case "user_id":
			conditions = append(conditions, "(chair_id = ? OR supervisor_id = ? OR assistant_supervisor_id = ? OR reviewer_id = ?)")
			for i := 0; i < 4; i++ {
				values = append(values, value[0])
			}
			continue
		case "supervisor_name":
			ids, err := m.GetSupervisorID(value[0])
			if err != nil {
				slog.Error(err.Error())
			}
			str := ""
			for _, id := range ids {
				str = str + "(supervisor_id = ?) OR"
				values = append(values, id)
			}
			if len(str) > 0 {
				str = str[0 : len(str)-3]
			}
			conditions = append(conditions, str)
			continue
		}
		if strings.Contains(key, "[") {
			field := key[:strings.Index(key, "[")]
			operator := key[strings.Index(key, "[")+1 : strings.Index(key, "]")]

			switch operator {
			case "gt":
				conditions = append(conditions, fmt.Sprintf("%s > ?", field))
			case "lt":
				conditions = append(conditions, fmt.Sprintf("%s < ?", field))
			case "gte":
				conditions = append(conditions, fmt.Sprintf("%s >= ?", field))
			case "lte":
				conditions = append(conditions, fmt.Sprintf("%s <= ?", field))
			}
			values = append(values, value[0])
		} else if strings.Contains(value[0], "|") {
			orValues := strings.Split(value[0], "|")
			placeholders := make([]string, len(orValues))
			for i, v := range orValues {
				placeholders[i] = "?"
				values = append(values, v)
			}
			conditions = append(conditions, fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ", ")))
		} else {
			conditions = append(conditions, fmt.Sprintf("%s = ?", key))
			values = append(values, value[0])
		}
	}
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}
	return baseQuery, values
}

func AddSQLOrder(baseQuery, sort_by string, desc_order bool) string {
	order := "DESC"
	if !desc_order {
		order = "ASC"
	}
	if sort_by == "" {
		sort_by = "id"
	}
	return baseQuery + fmt.Sprintf(" ORDER BY %v %v", sort_by, order)
}
