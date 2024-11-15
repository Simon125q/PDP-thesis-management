package types

type RealizedThesis struct {
	Id                    int
	ThesisNumber          string
	ExamDate              string
	AverageStudyGrade     string
	CompetencyExamGrade   string
	DiplomaExamGrade      string
	FinalStudyResult      string
	FinalStudyResultText  string
	ThesisTitlePolish     string
	ThesisTitleEnglish    string
	ThesisLanguage        string
	Library               string
	StudentId             int
	ChairId               int
	SupervisorId          int
	AssistantSupervisorId int
	ReviewerId            int
	HourlySettlementId    int
}
