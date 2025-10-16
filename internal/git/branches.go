package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// BranchInfo represents information about a Git branch
type BranchInfo struct {
	Name      string
	IsCurrent bool
	IsRemote  bool
}

// GetBranchesWithInfo returns detailed information about all branches
func (g *LocalGitRepository) GetBranchesWithInfo() ([]BranchInfo, error) {
	if !g.IsGitRepository() {
		return nil, GitError{
			Operation: "branch",
			Message:   "not a git repository",
		}
	}

	cmd := exec.Command("git", "branch", "-a", "--format=%(refname:short)|%(HEAD)")
	output, err := cmd.Output()
	if err != nil {
		return nil, GitError{
			Operation: "branch",
			Message:   "failed to list branches with info: " + err.Error(),
		}
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var branches []BranchInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}

		branchName := parts[0]
		isCurrent := parts[1] == "*"
		isRemote := strings.HasPrefix(branchName, "origin/") || strings.Contains(branchName, "/")

		// Skip remote tracking branches for local branch listing
		if !isRemote {
			branches = append(branches, BranchInfo{
				Name:      branchName,
				IsCurrent: isCurrent,
				IsRemote:  isRemote,
			})
		}
	}

	return branches, nil
}

// BranchSearchResult represents the result of a branch search operation
type BranchSearchResult struct {
	Branches []string
	HasResults bool
	SearchTerm string
}

// FilterBranches filters branches based on a search term using fuzzy matching
func FilterBranches(branches []string, searchTerm string) BranchSearchResult {
	result := BranchSearchResult{
		SearchTerm: searchTerm,
		HasResults: true,
	}

	// If no search term, return all branches
	if searchTerm == "" {
		result.Branches = branches
		return result
	}

	searchTerm = strings.ToLower(searchTerm)
	var filtered []string

	// Simple fuzzy matching - check if all characters of search term appear in order
	for _, branch := range branches {
		if fuzzyMatch(strings.ToLower(branch), searchTerm) {
			filtered = append(filtered, branch)
		}
	}

	result.Branches = filtered
	result.HasResults = len(filtered) > 0

	return result
}

// fuzzyMatch performs fuzzy matching between a branch name and search term
// Returns true if all characters in searchTerm appear in branchName in order
func fuzzyMatch(branchName, searchTerm string) bool {
	if searchTerm == "" {
		return true
	}

	branchRunes := []rune(branchName)
	searchRunes := []rune(searchTerm)
	
	branchIndex := 0
	searchIndex := 0

	for branchIndex < len(branchRunes) && searchIndex < len(searchRunes) {
		if branchRunes[branchIndex] == searchRunes[searchIndex] {
			searchIndex++
		}
		branchIndex++
	}

	return searchIndex == len(searchRunes)
}

// FilterBranchesRealtime provides real-time filtering with enhanced fuzzy matching
func FilterBranchesRealtime(branches []string, searchTerm string) BranchSearchResult {
	result := BranchSearchResult{
		SearchTerm: searchTerm,
		HasResults: true,
	}

	// If no search term, return all branches
	if searchTerm == "" {
		result.Branches = branches
		return result
	}

	searchTerm = strings.ToLower(searchTerm)
	var filtered []string

	// Enhanced matching: exact substring match gets priority, then fuzzy match
	var exactMatches []string
	var fuzzyMatches []string

	for _, branch := range branches {
		lowerBranch := strings.ToLower(branch)
		
		// Check for exact substring match first
		if strings.Contains(lowerBranch, searchTerm) {
			exactMatches = append(exactMatches, branch)
		} else if fuzzyMatch(lowerBranch, searchTerm) {
			// Only add to fuzzy matches if not already in exact matches
			fuzzyMatches = append(fuzzyMatches, branch)
		}
	}

	// Combine results with exact matches first
	filtered = append(exactMatches, fuzzyMatches...)

	result.Branches = filtered
	result.HasResults = len(filtered) > 0

	return result
}

// GetEmptySearchMessage returns an appropriate message when no branches match the search
func (result BranchSearchResult) GetEmptySearchMessage() string {
	if result.SearchTerm == "" {
		return "No branches available"
	}
	return "No branches found matching '" + result.SearchTerm + "'"
}

// GetSearchSummary returns a summary of the search results
func (result BranchSearchResult) GetSearchSummary() string {
	if result.SearchTerm == "" {
		return ""
	}
	
	count := len(result.Branches)
	if count == 0 {
		return result.GetEmptySearchMessage()
	} else if count == 1 {
		return "1 branch found"
	} else {
		return fmt.Sprintf("%d branches found", count)
	}
}