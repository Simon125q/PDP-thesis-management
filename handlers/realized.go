package handlers

import (
	"net/http"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/views/realized"
)

func HandleRealized(w http.ResponseWriter, r *http.Request) error {
	emp_data, err := server.MyS.DB.AllUniversityEmployee()
	if err != nil {
		return err
	}
	return Render(w, r, realized.Index(emp_data))
}
