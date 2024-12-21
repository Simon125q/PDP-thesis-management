package types

type Specialization struct {
	Id   int
	Name string
}

type SpecializationErrors struct {
	Id      int
	Name    string
	Correct bool
}

type SpecializationEntry struct {
	Id   int
	Name string
}

type SpecializationEntryErrors struct {
	Id      int
	Name    string
	Correct bool
}
