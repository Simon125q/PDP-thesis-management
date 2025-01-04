package types

type OngoingThesis struct {
	Id                               int
	ThesisNumber                     string
	ThesisTitlePolish                string
	ThesisTitleEnglish               string
	ThesisLanguage                   string
	StudentId                        int
	SupervisorAcademicTitle          string
	SupervisorId                     int
	AssistantSupervisorAcademicTitle string
	AssistantSupervisorId            int
}

type OngoingThesisEntry struct {
	Id                               int
	ThesisNumber                     string
	ThesisTitlePolish                string
	ThesisTitleEnglish               string
	ThesisLanguage                   string
	Student                          Student
	SupervisorAcademicTitle          string
	Supervisor                       UniversityEmployeeEntry
	AssistantSupervisorAcademicTitle string
	AssistantSupervisor              UniversityEmployeeEntry
}

type OngoingThesisEntryErrors struct {
	Id                               int
	ThesisNumber                     string
	ThesisTitlePolish                string
	ThesisTitleEnglish               string
	ThesisLanguage                   string
	Student                          StudentErrors
	SupervisorAcademicTitle          string
	Supervisor                       UniversityEmployeeErrors
	AssistantSupervisorAcademicTitle string
	AssistantSupervisor              UniversityEmployeeErrors
	Correct                          bool
	InternalError                    bool
}
