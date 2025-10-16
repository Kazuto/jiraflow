package jira

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"jiraflow/internal/errors"
)

// JiraClient interface defines Jira operations
type JiraClient interface {
	GetTicketTitle(ticketID string) (string, error)
	IsAvailable() bool
}

// JiraError is an alias for the centralized JiraError type
type JiraError = errors.JiraError

// CLIClient implements JiraClient using the Jira CLI
type CLIClient struct{}

// NewCLIClient creates a new Jira CLI client
func NewCLIClient() *CLIClient {
	return &CLIClient{}
}

// IsAvailable checks if the Jira CLI is installed and available
func (c *CLIClient) IsAvailable() bool {
	_, err := exec.LookPath("jira")
	return err == nil
}

// GetTicketTitle fetches the ticket title using the Jira CLI
func (c *CLIClient) GetTicketTitle(ticketID string) (string, error) {
	if !c.IsAvailable() {
		return "", errors.NewJiraError(ticketID, "jira CLI not found - please install jira CLI or provide title manually", true)
	}

	// Execute jira issue view command with raw JSON output
	cmd := exec.Command("jira", "issue", "view", ticketID, "--raw")
	output, err := cmd.Output()
	if err != nil {
		// Try to get more specific error information
		if exitError, ok := err.(*exec.ExitError); ok {
			stderr := string(exitError.Stderr)
			if strings.Contains(stderr, "not found") || strings.Contains(stderr, "does not exist") {
				return "", errors.NewJiraError(ticketID, fmt.Sprintf("ticket %s not found", ticketID), true)
			}
			if strings.Contains(stderr, "authentication") || strings.Contains(stderr, "unauthorized") {
				return "", errors.NewJiraError(ticketID, "authentication failed - please run 'jira init' to configure credentials", true)
			}
			return "", errors.NewJiraError(ticketID, fmt.Sprintf("failed to fetch ticket: %s", stderr), true)
		}
		return "", errors.NewJiraError(ticketID, fmt.Sprintf("failed to execute jira command: %v", err), true)
	}

	// Parse the JSON output to extract the title
	title, err := c.parseJSONTitle(string(output))
	if err != nil {
		return "", errors.NewJiraError(ticketID, fmt.Sprintf("failed to parse ticket title: %v", err), true)
	}

	if title == "" {
		return "", errors.NewJiraError(ticketID, "ticket title is empty", true)
	}

	return title, nil
}



// parseJSONTitle parses the JSON output from jira CLI --raw command
func (c *CLIClient) parseJSONTitle(output string) (string, error) {
	// Parse the JSON response from jira issue view --raw
	var ticket struct {
		Fields struct {
			Summary string `json:"summary"`
		} `json:"fields"`
	}
	
	if err := json.Unmarshal([]byte(output), &ticket); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Return the summary field which contains the ticket title
	return ticket.Fields.Summary, nil
}

// MockClient implements JiraClient for testing purposes
type MockClient struct {
	Available bool
	Tickets   map[string]string
	Error     error
}

// NewMockClient creates a new mock Jira client
func NewMockClient() *MockClient {
	return &MockClient{
		Available: true,
		Tickets:   make(map[string]string),
	}
}

// IsAvailable returns the mock availability status
func (m *MockClient) IsAvailable() bool {
	return m.Available
}

// GetTicketTitle returns the mock ticket title or error
func (m *MockClient) GetTicketTitle(ticketID string) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}
	
	if !m.Available {
		return "", errors.NewJiraError(ticketID, "jira CLI not available", true)
	}
	
	title, exists := m.Tickets[ticketID]
	if !exists {
		return "", errors.NewJiraError(ticketID, "ticket not found", true)
	}
	
	return title, nil
}

// SetTicket adds a ticket to the mock client
func (m *MockClient) SetTicket(ticketID, title string) {
	m.Tickets[ticketID] = title
}

// SetError sets an error to be returned by GetTicketTitle
func (m *MockClient) SetError(err error) {
	m.Error = err
}

// SetAvailable sets the availability status
func (m *MockClient) SetAvailable(available bool) {
	m.Available = available
}