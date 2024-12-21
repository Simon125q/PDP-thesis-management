package validators

import "thesis-management-app/types"

func ValidateCourse(course types.Course) (types.CourseErrors, bool) {
	ok := true
	fErr, fOk := ValidateCourseName(course.Name)

	if !fOk {
		ok = false
	}
	return types.CourseErrors{
		Name: fErr,
	}, ok
}

func ValidateCourseName(courseName string) (string, bool) {
	return "", true
	//me love some good checking
}
