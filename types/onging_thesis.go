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
	Note                             Note
}

type OngoingThesisEntryErrors struct {
	Id                               int
	ThesisNumber                     string
	ThesisTitlePolish                string
	ThesisTitleEnglish               string
	ThesisLanguage                   string
	Student                          StudentErrors
	SupervisorAcademicTitle          string
	Supervisor                       UniversityEmployeeEntryErrors
	AssistantSupervisorAcademicTitle string
	AssistantSupervisor              UniversityEmployeeEntryErrors
	Correct                          bool
	InternalError                    bool
}
