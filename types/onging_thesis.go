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
