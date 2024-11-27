package types

type HourlySettlement struct {
	Id                       int
	SupervisorHours          float64
	AssistantSupervisorHours float64
	ReviewerHours            float64
}

type HourlySettlementErrors struct {
	SupervisorHours          string
	AssistantSupervisorHours string
	ReviewerHours            string
}
