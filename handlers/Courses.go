package handlers

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/pkgs/validators"
	"thesis-management-app/types"
	"thesis-management-app/views/settings"
)

func HandleCourses(w http.ResponseWriter, r *http.Request) error {
	courses, err := server.MyS.DB.AllCourses()
	if err != nil {
		slog.Error("Error fetching courses", "err", err)
		return err
	}
	return Render(w, r, settings.ResultsCourses(courses))
}

func HandleCoursesEntry(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HandleCoursesEntry", "id_param", id_param)

	course, err := server.MyS.DB.CourseById(id_param)
	if err != nil {
		slog.Error("Error fetching course", "err", err)
		return err
	}

	return Render(w, r, settings.Entry_Course(course))
}

func HandleCoursesNew(w http.ResponseWriter, r *http.Request) error {
	course := *extractCourseFromForm(r)

	errors, ok := validators.ValidateCourse(course)
	if !ok {
		errors.Correct = false
		return Render(w, r, settings.NewEntrySwap_Course(types.Course{}, course, errors))
	}

	slog.Info("add course", "correct", true)

	courseId, err := server.MyS.DB.InsertCourse(course)
	if err != nil {
		slog.Error("course to db", "err", err)
		errors.Correct = false
		return Render(w, r, settings.NewEntrySwap_Course(types.Course{}, course, errors))
	}

	slog.Info("course to db", "new_id", courseId)
	course.Id = int(courseId)
	errors.Correct = true

	return Render(w, r, settings.NewEntrySwap_Course(course, types.Course{}, errors))
}

func extractCourseFromForm(r *http.Request) *types.Course {
	return &types.Course{
		Name: r.FormValue("name"),
	}
}

func HandleCoursesDetails(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HSDetails", "id_param", id_param)
	course_data, err := server.MyS.DB.CourseById(id_param)
	slog.Info("quere", "q", r.URL.Query())
	if err != nil {
		return err
	}
	slog.Info("HSettingsCourseDetails", "course", course_data)
	return Render(w, r, settings.Details_Course(course_data, types.CourseErrors{}))
}

func HandleCoursesUpdate(w http.ResponseWriter, r *http.Request) error {
	id_param := chi.URLParam(r, "id")
	slog.Info("HandleCoursesUpdate", "id_param", id_param)

	course := types.Course{
		Name: r.FormValue("name"),
	}

	errors, ok := validators.ValidateCourse(course)
	if !ok {
		errors.Correct = false
		return Render(w, r, settings.Entry_Course(course))
	}

	var err error
	course.Id, err = strconv.Atoi(id_param)
	if err != nil {
		slog.Error("Error parsing course ID", "err", err)
		return Render(w, r, settings.Entry_Course(course))
	}

	err = server.MyS.DB.UpdateCourse(course) // tu było &courses ale cos mu nie gralo, moze tu jakis blad
	if err != nil {
		slog.Error("Error updating course", "err", err)
		return Render(w, r, settings.Entry_Course(course))
	}

	slog.Info("Course updated", "courseId", course.Id)
	errors.Correct = true
	return Render(w, r, settings.Entry_Course(course))
}

func HandleCoursesClearNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, settings.EmptySpace_Course())
}

func HandleCoursesGetNew(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, settings.NewEntry_Course(types.Course{}, types.CourseErrors{}))
}
