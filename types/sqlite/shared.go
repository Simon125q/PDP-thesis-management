package sqlite

import (
	"fmt"
	"net/url"
	"strings"
)

func AddSQLQueryParameters(baseQuery string, params url.Values) (string, []interface{}) {
	var conditions []string
	var values []interface{}

	for key, value := range params {
		if key == "user_id" {
			conditions = append(conditions, "(chair_id = ? OR supervisor_id = ? OR assistant_supervisor_id = ? OR reviewer_id = ?)")
			for i := 0; i < 4; i++ {
				values = append(values, value[0])
			}
			continue
		}
		if strings.Contains(key, "[") {
			field := key[:strings.Index(key, "]")]
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
