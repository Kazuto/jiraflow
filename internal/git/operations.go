package git

import (
	"os/exec"
	"strings"

	"jiraflow/internal/errors"
)

// GitRepository interface defines Git operations
type GitRepository interface {
	GetLocalBranches() ([]string, error)
	GetBranchesWithInfo() ([]BranchInfo, error)
	GetCurrentBranch() (string, error)
	CreateBranch(name, baseBranch string) error
	CheckoutBranch(name string) error
	IsGitRepository() bool
	SearchBranches(searchTerm string) (BranchSearchResult, error)
}

// GitError is an alias for the centralized GitError type
type GitError = errors.GitError

// LocalGitRepository implements GitRepository for local Git operations
type LocalGitRepository struct{}

// NewLocalGitRepository creates a new LocalGitRepository instance
func NewLocalGitRepository() *LocalGitRepository {
	return &LocalGitRepository{}
}

// IsGitRepository checks if the current directory is a Git repository
func (g *LocalGitRepository) IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// GetLocalBranches returns a list of local Git branches
func (g *LocalGitRepository) GetLocalBranches() ([]string, error) {
	if !g.IsGitRepository() {
		return nil, errors.NewGitError("branch", "not a git repository", false)
	}

	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.NewGitError("branch", "failed to list branches: "+err.Error(), true)
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	var localBranches []string
	
	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		if branch != "" {
			localBranches = append(localBranches, branch)
		}
	}

	return localBranches, nil
}

// GetCurrentBranch returns the name of the current Git branch
func (g *LocalGitRepository) GetCurrentBranch() (string, error) {
	if !g.IsGitRepository() {
		return "", errors.NewGitError("branch", "not a git repository", false)
	}

	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.NewGitError("branch", "failed to get current branch: "+err.Error(), true)
	}

	currentBranch := strings.TrimSpace(string(output))
	if currentBranch == "" {
		return "", errors.NewGitError("branch", "no current branch (detached HEAD?)", true)
	}

	return currentBranch, nil
}

// CreateBranch creates a new Git branch from the specified base branch and checks it out
func (g *LocalGitRepository) CreateBranch(name, baseBranch string) error {
	if !g.IsGitRepository() {
		return errors.NewGitError("branch", "not a git repository", false)
	}

	if name == "" {
		return errors.NewGitError("branch", "branch name cannot be empty", false)
	}

	if baseBranch == "" {
		return errors.NewGitError("branch", "base branch cannot be empty", false)
	}

	// Create and checkout the branch in one command
	cmd := exec.Command("git", "checkout", "-b", name, baseBranch)
	if err := cmd.Run(); err != nil {
		return errors.NewGitError("branch", "failed to create and checkout branch '"+name+"' from '"+baseBranch+"': "+err.Error(), true)
	}

	return nil
}

// CheckoutBranch switches to the specified Git branch
func (g *LocalGitRepository) CheckoutBranch(name string) error {
	if !g.IsGitRepository() {
		return errors.NewGitError("checkout", "not a git repository", false)
	}

	if name == "" {
		return errors.NewGitError("checkout", "branch name cannot be empty", false)
	}

	cmd := exec.Command("git", "checkout", name)
	if err := cmd.Run(); err != nil {
		return errors.NewGitError("checkout", "failed to checkout branch '"+name+"': "+err.Error(), true)
	}

	return nil
}

// SearchBranches searches for branches matching the given search term
func (g *LocalGitRepository) SearchBranches(searchTerm string) (BranchSearchResult, error) {
	branches, err := g.GetLocalBranches()
	if err != nil {
		return BranchSearchResult{}, err
	}

	result := FilterBranchesRealtime(branches, searchTerm)
	return result, nil
}