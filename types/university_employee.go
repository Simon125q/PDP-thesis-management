package types

type UniversityEmployee struct {
	Id                   int
	FirstName            string
	LastName             string
	CurrentAcademicTitle string
	DepartmentUnit       string
}

type UniversityEmployeeErrors struct {
	FirstName            string
	LastName             string
	CurrentAcademicTitle string
	DepartmentUnit       string
	Correct              bool
}

type UniversityEmployeeEntry struct {
	Id                   int
	FirstName            string
	LastName             string
	CurrentAcademicTitle string
	DepartmentUnit       string
}

type UniversityEmployeeEntryErrors struct {
	Id                   int
	FirstName            string
	LastName             string
	CurrentAcademicTitle string
	DepartmentUnit       string
	Correct              bool
}
