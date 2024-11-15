package handlers

import (
	"net/http"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/views/realized"
)

func HandleRealized(w http.ResponseWriter, r *http.Request) error {
	thes_data, err := server.MyS.DB.AllRealizedThesis("id", false)
	if err != nil {
		return err
	}
	return Render(w, r, realized.Index(thes_data))
}
