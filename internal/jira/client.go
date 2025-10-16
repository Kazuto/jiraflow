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

	// Execute jira view command with JSON output
	cmd := exec.Command("jira", "view", ticketID, "--plain")
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

	// Parse the output to extract the title
	title, err := c.parseTicketTitle(string(output))
	if err != nil {
		return "", errors.NewJiraError(ticketID, fmt.Sprintf("failed to parse ticket title: %v", err), true)
	}

	if title == "" {
		return "", errors.NewJiraError(ticketID, "ticket title is empty", true)
	}

	return title, nil
}

// parseTicketTitle extracts the title from jira CLI output
func (c *CLIClient) parseTicketTitle(output string) (string, error) {
	lines := strings.Split(output, "\n")
	
	// Look for the summary/title line in the plain output
	// The jira CLI plain output typically has the format:
	// KEY: TICKET-123
	// Summary: The ticket title
	// Status: In Progress
	// ...
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "summary:") {
			// Extract the title after "Summary:"
			title := strings.TrimSpace(strings.TrimPrefix(line, "Summary:"))
			title = strings.TrimSpace(strings.TrimPrefix(title, "summary:"))
			return title, nil
		}
		// Sometimes the title might be on the next line
		if strings.ToLower(line) == "summary:" && i+1 < len(lines) {
			title := strings.TrimSpace(lines[i+1])
			return title, nil
		}
	}

	// Fallback: try JSON parsing if plain format doesn't work
	return c.parseJSONTitle(output)
}

// parseJSONTitle attempts to parse JSON output as fallback
func (c *CLIClient) parseJSONTitle(output string) (string, error) {
	// Try to parse as JSON in case the output format is different
	var ticket struct {
		Fields struct {
			Summary string `json:"summary"`
		} `json:"fields"`
	}
	
	if err := json.Unmarshal([]byte(output), &ticket); err == nil {
		return ticket.Fields.Summary, nil
	}

	// If JSON parsing fails, return error
	return "", fmt.Errorf("could not parse ticket title from output")
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