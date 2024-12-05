package types

type Student struct {
	Id             int
	StudentNumber  string
	FirstName      string
	LastName       string
	FieldOfStudy   string
	Degree         string
	Specialization string
	ModeOfStudies  string
}

type StudentErrors struct {
	StudentNumber  string
	FirstName      string
	LastName       string
	FieldOfStudy   string
	Degree         string
	Specialization string
	ModeOfStudies  string
}
