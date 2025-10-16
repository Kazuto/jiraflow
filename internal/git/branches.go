package git

import (
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

// FilterBranches filters branches based on a search term using fuzzy matching
func FilterBranches(branches []string, searchTerm string) []string {
	if searchTerm == "" {
		return branches
	}

	searchTerm = strings.ToLower(searchTerm)
	var filtered []string

	for _, branch := range branches {
		if strings.Contains(strings.ToLower(branch), searchTerm) {
			filtered = append(filtered, branch)
		}
	}

	return filtered
}