package types

type Course struct {
	Id   int
	Name string
}

type CourseErrors struct {
	Id      int
	Name    string
	Correct bool
}

type CourseEntry struct {
	Id   int
	Name string
}

type CourseEntryErrors struct {
	Id      int
	Name    string
	Correct bool
}
