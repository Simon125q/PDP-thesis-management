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
	Archived                         string
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
	Archived                         string
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
	Checklist                        string
	Correct                          bool
	InternalError                    bool
}
