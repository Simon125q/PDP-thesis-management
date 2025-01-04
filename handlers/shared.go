package handlers

import (
	"log/slog"
	"math"
	"net/http"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/types"

	"github.com/a-h/templ"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request) error

func Make(h HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("HTTP handler error", "err", err, "path", r.URL.Path)
		}
	}
}

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	return c.Render(r.Context(), w)
}

func hxRedirect(w http.ResponseWriter, r *http.Request, to string) error {
	if len(r.Header.Get("HX-Request")) > 0 {
		w.Header().Set("HX-Redirect", to)
		w.WriteHeader(http.StatusSeeOther)
		return nil
	}
	http.Redirect(w, r, to, http.StatusSeeOther)
	return nil
}

func getAutehenticatedUser(r *http.Request) types.AuthenticatedUser {
	user, ok := r.Context().Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}
	return user
}

func getEmployeeId(emp types.UniversityEmployeeEntry) (int, error) {
	empId, err := server.MyS.DB.EmployeeIdByName(emp.FirstName + " " + emp.LastName)
	slog.Info("getEmployeeId", "empName", emp.FirstName+" "+emp.LastName)
	slog.Info("getEmployeeId", "empId", empId)
	if err != nil {
		slog.Error("getEmployeeId", "err", err)
	}
	if empId == 0 {
		if emp.FirstName != "" && emp.LastName != "" {
			var id int64
			id, err = server.MyS.DB.InsertUniversityEmployee(emp)
			if err != nil {
				slog.Error("getEmployeeId", "err", err)
			}
			slog.Info("getEmployeeId", "inserting new emp, id", id)
			empId = int(id)
		}
	}
	return empId, err
}

func paginate[K any](data []K, page, pageSize int) ([]K, int) {
	slog.Info("paginate", "page", page)
	slog.Info("paginate", "len(data)", len(data))
	totalItems := len(data)
	if totalItems <= 0 {
		return []K{}, 0
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if end > totalItems {
		end = totalItems
	}
	return data[start:end], totalPages
}
