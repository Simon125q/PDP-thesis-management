package types

type RealizedThesis struct {
	Id                               int
	ThesisNumber                     string
	ExamDate                         string
	ExamTime                         string
	AverageStudyGrade                string
	CompetencyExamGrade              string
	DiplomaExamGrade                 string
	FinalStudyResult                 string
	FinalStudyResultText             string
	ThesisTitlePolish                string
	ThesisTitleEnglish               string
	ThesisLanguage                   string
	Library                          string
	StudentId                        int
	ChairAcademicTitle               string
	ChairId                          int
	SupervisorAcademicTitle          string
	SupervisorId                     int
	AssistantSupervisorAcademicTitle string
	AssistantSupervisorId            int
	ReviewerAcademicTitle            string
	ReviewerId                       int
	HourlySettlementId               int
}

type RealizedThesisEntry struct {
	Id                               int
	ThesisNumber                     string
	ExamDate                         string
	ExamTime                         string
	AverageStudyGrade                string
	CompetencyExamGrade              string
	DiplomaExamGrade                 string
	FinalStudyResult                 string
	FinalStudyResultText             string
	ThesisTitlePolish                string
	ThesisTitleEnglish               string
	ThesisLanguage                   string
	Library                          string
	Student                          Student
	ChairAcademicTitle               string
	Chair                            UniversityEmployeeEntry
	SupervisorAcademicTitle          string
	Supervisor                       UniversityEmployeeEntry
	AssistantSupervisorAcademicTitle string
	AssistantSupervisor              UniversityEmployeeEntry
	ReviewerAcademicTitle            string
	Reviewer                         UniversityEmployeeEntry
	HourlySettlement                 HourlySettlement
	Note                             Note
}

type RealizedThesisEntryErrors struct {
	ThesisNumber                     string
	ExamDate                         string
	ExamTime                         string
	AverageStudyGrade                string
	CompetencyExamGrade              string
	DiplomaExamGrade                 string
	FinalStudyResult                 string
	FinalStudyResultText             string
	ThesisTitlePolish                string
	ThesisTitleEnglish               string
	ThesisLanguage                   string
	Library                          string
	Student                          StudentErrors
	ChairAcademicTitle               string
	Chair                            UniversityEmployeeEntryErrors
	SupervisorAcademicTitle          string
	Supervisor                       UniversityEmployeeEntryErrors
	AssistantSupervisorAcademicTitle string
	AssistantSupervisor              UniversityEmployeeEntryErrors
	ReviewerAcademicTitle            string
	Reviewer                         UniversityEmployeeEntryErrors
	HourlySettlement                 HourlySettlementErrors
	Correct                          bool
	InternalError                    bool
}
