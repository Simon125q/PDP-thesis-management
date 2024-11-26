package types

type RealizedThesis struct {
	Id                               int
	ThesisNumber                     string
	ExamDate                         string
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
	Chair                            UniversityEmployee
	SupervisorAcademicTitle          string
	Supervisor                       UniversityEmployee
	AssistantSupervisorAcademicTitle string
	AssistantSupervisor              UniversityEmployee
	ReviewerAcademicTitle            string
	Reviewer                         UniversityEmployee
	HourlySettlement                 HourlySettlement
}
