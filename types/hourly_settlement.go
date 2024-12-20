package types

type HourlySettlement struct {
	Id                              int
	SupervisorHours                 int
	SupervisorHoursSettled          int
	AssistantSupervisorHours        int
	AssistantSupervisorHoursSettled int
	ReviewerHours                   int
	ReviewerHoursSettled            int
}

type HourlySettlementErrors struct {
	SupervisorHours          string
	AssistantSupervisorHours string
	ReviewerHours            string
	Total                    string
}
