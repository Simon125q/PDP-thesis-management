package types

type UniversityEmployee struct {
	Id                   int
	FirstName            string
	LastName             string
	CurrentAcademicTitle string
	DepartmentUnit       string
	ThesisCount          string
}

type UniversityEmployeeErrors struct {
	FirstName            string
	LastName             string
	CurrentAcademicTitle string
	DepartmentUnit       string
	ThesisCount          string
	Correct              bool
}

type UniversityEmployeeEntry struct {
	Id                   int
	FirstName            string
	LastName             string
	CurrentAcademicTitle string
	DepartmentUnit       string
	ThesisCount          string
}

type UniversityEmployeeEntryErrors struct {
	Id                   int
	FirstName            string
	LastName             string
	CurrentAcademicTitle string
	DepartmentUnit       string
	ThesisCount          string
	Correct              bool
}
