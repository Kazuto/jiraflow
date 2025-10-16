package jira

// JiraClient interface defines Jira operations
type JiraClient interface {
	GetTicketTitle(ticketID string) (string, error)
	IsAvailable() bool
}

// JiraError represents a Jira operation error
type JiraError struct {
	TicketID string
	Message  string
}

func (e JiraError) Error() string {
	return "jira error for ticket " + e.TicketID + ": " + e.Message
}