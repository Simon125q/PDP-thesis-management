package sqlite

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

func (m *Model) GetConditionValuesForTasks(conditions []string, values []interface{}) ([]string, []interface{}) {
	queryString := fmt.Sprintf("SELECT id FROM Thesis_To_Be_Completed")
	slog.Info(queryString)
	ids, err := m.DB.Query(queryString)
	defer ids.Close()
	if err != nil {
		slog.Error(err.Error())
	}
	var maxID int
	for ids.Next() {
		if ids != nil {
			var id int
			ids.Scan(&id)
			if id > maxID {
				maxID = id
			}
		}
	}
	var ongoingIDTable []bool
	for range maxID + 1 {
		ongoingIDTable = append(ongoingIDTable, false)
	}
	for i := range maxID + 1 {
		if i == 0 {
			continue
		}
		queryStr := fmt.Sprintf("SELECT is_completed FROM Task WHERE thesis_to_be_completed_id = %v", i)
		slog.Info(queryStr)
		val, qErr := m.DB.Query(queryStr)
		defer val.Close()
		if qErr != nil {
			slog.Error(qErr.Error())
		}
		allTicked := true
		hasTasks := false
		for val.Next() {
			hasTasks = true
			var completed int
			val.Scan(&completed)
			if completed != 1 {
				allTicked = false
			}
		}
		if !hasTasks {
			allTicked = true
		}
		if !allTicked {
			ongoingIDTable[i] = true
		}
	}
	str := ""
	for i := range ongoingIDTable {
		if ongoingIDTable[i] {
			str = str + " thesis_id = ? OR"
			values = append(values, i)
		}
	}
	if len(str) > 0 {
		str = str[0 : len(str)-3]
		str = "(" + str + ")"
	}
	conditions = append(conditions, str)
	//slog.Info("TASK CONDS AND VALUES", "conds", conditions)
	//slog.Info("TASK CONDS AND VALUES", "vals", values)
	return conditions, values
}

func (m *Model) GetStudentID(value string, column string) ([]string, error) {
	queryString := fmt.Sprintf("SELECT id FROM 'Student' WHERE %v LIKE '%v'", column, value)
	slog.Info(queryString)
	rows, err := m.DB.Query(queryString)
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

func (m *Model) GetPersonID(name string, personRank string) ([]string, error) {
	name = "%" + strings.Replace(name, " ", "%", -1) + "%"
	queryString := fmt.Sprintf("SELECT id FROM '%v' WHERE CONCAT(first_name, ' ', last_name) LIKE '%v'", personRank, name)
	slog.Info(queryString)
	rows, err := m.DB.Query(queryString)
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
	if len(ids) == 0 {
		ids = append(ids, "0")
	}
	return ids, nil
}

func (m *Model) GetConditionValuesFromName(name string, personRank string, column string, conditions []string, values []interface{}) ([]string, []interface{}) {
	slog.Info("GetConditionValuesFromName", "personRank", personRank)
	ids, err := m.GetPersonID(name, personRank)
	if err != nil {
		slog.Error("GetConditionValuesFromName", "err", err)
		return nil, nil
	}
	slog.Info("GetConditionValuesFromName", "ids", ids)
	str := "("
	for _, id := range ids {
		str = str + " " + column + " = ? OR" //
		values = append(values, id)
	}
	if len(str) > 1 {
		str = str[0 : len(str)-3]
		str = str + ")"
		conditions = append(conditions, str)
	}
	slog.Info("GetConditionValuesFromName", "conditions", conditions)
	slog.Info("GetConditionValuesFromName", "vals", values)
	return conditions, values
}

func (m *Model) GetConditionValuesFromStudent(value string, valueType string, column string, conditions []string, values []interface{}) ([]string, []interface{}) {
	ids, err := m.GetStudentID(value, valueType)
	if err != nil {
		slog.Error(err.Error())
	}
	str := "("
	if len(ids) <= 0 {
		str = str + " " + column + " = 'BLANK' OR"
	}
	for _, id := range ids {
		str = str + " " + column + " = ? OR"
		values = append(values, id)
	}
	if len(str) > 1 {
		str = str[0 : len(str)-3]
		str = str + ")"
		conditions = append(conditions, str)
	}
	return conditions, values
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
		case "ongoing_user_id":
			conditions = append(conditions, "(supervisor_id = ? OR assistant_supervisor_id = ?)")
			for i := 0; i < 2; i++ {
				values = append(values, value[0])
			}
			continue
		case "mean-grade-min":
			conditionStr := fmt.Sprintf("average_study_grade >= %v", value[0])
			conditions = append(conditions, conditionStr)
			continue
		case "mean-grade-max":
			conditionStr := fmt.Sprintf("average_study_grade <= %v", value[0])
			conditions = append(conditions, conditionStr)
			continue
		case "thesis_title":
			thesisTitle := value[0]
			thesisTitle = "%" + thesisTitle + "%"
			conditionStr := fmt.Sprintf("(thesis_title_polish LIKE '%v' OR thesis_title_english LIKE '%v')", thesisTitle, thesisTitle)
			conditions = append(conditions, conditionStr)
			continue
		case "student_name":
			conditions, values = m.GetConditionValuesFromName(value[0], "Student", "student_id", conditions, values)
			continue
		case "student_number":
			conditions, values = m.GetConditionValuesFromStudent(value[0]+"%", "student_number", "student_id", conditions, values)
			continue
		case "supervisor_name":
			conditions, values = m.GetConditionValuesFromName(value[0], "University_Employee", "supervisor_id", conditions, values)
			continue
		case "assistant_supervisor_name":
			conditions, values = m.GetConditionValuesFromName(value[0], "University_Employee", "assistant_supervisor_id", conditions, values)
			continue
		case "reviewer_name":
			conditions, values = m.GetConditionValuesFromName(value[0], "University_Employee", "reviewer_id", conditions, values)
			continue
		case "course":
			conditions, values = m.GetConditionValuesFromStudent(value[0], "field_of_study", "student_id", conditions, values)
			continue
		case "mode_of_studies":
			conditions, values = m.GetConditionValuesFromStudent(value[0], "mode_of_studies", "student_id", conditions, values)
			continue
		case "degree":
			conditions, values = m.GetConditionValuesFromStudent(value[0], "degree", "student_id", conditions, values)
			continue
		case "archived_ongoing":
			switch value[0] {
			case "archived":
				conditions = append(conditions, "(topic_scan = 'true')")
				continue
			case "notarchived":
				conditions = append(conditions, "(topic_scan = 'false' OR topic_scan IS NULL)")
				continue
			}
			continue
		case "are_hours_settled": // retired
			conditions = append(conditions, "(is_supervisor_settled = ? OR is_assistant_supervisor_settled = ? OR is_reviewer_settled = ?)")
			for i := 0; i < 3; i++ {
				values = append(values, "0")
			}
			continue
		case "are_hours_settled_supervisor":
			conditions = append(conditions, "(is_supervisor_settled = ?)")
			values = append(values, "0")
			continue
		case "are_hours_settled_assistant":
			conditions = append(conditions, "(is_assistant_supervisor_settled = ?)")
			values = append(values, "0")
			continue
		case "are_hours_settled_reviewer":
			conditions = append(conditions, "(is_reviewer_settled = ?)")
			values = append(values, "0")
			continue
		case "are_tasks_settled":
			conditions, values = m.GetConditionValuesForTasks(conditions, values)
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
		slog.Info("adding conds to query", "conds", conditions)
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}
	return baseQuery, values
}
func AddSQLPagination(baseQuery string, page_num, page_limit int) string {
	if page_num == -1 {
		return baseQuery
	}
	offset := (page_num - 1) * page_limit
	return baseQuery + fmt.Sprintf(" LIMIT %v OFFSET %v", page_limit, offset)
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
