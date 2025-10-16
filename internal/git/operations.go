package git

import (
	"os/exec"
	"strings"
)

// GitRepository interface defines Git operations
type GitRepository interface {
	GetLocalBranches() ([]string, error)
	GetCurrentBranch() (string, error)
	CreateBranch(name, baseBranch string) error
	CheckoutBranch(name string) error
	IsGitRepository() bool
}

// GitError represents a Git operation error
type GitError struct {
	Operation string
	Message   string
}

func (e GitError) Error() string {
	return "git " + e.Operation + ": " + e.Message
}

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
		return nil, GitError{
			Operation: "branch",
			Message:   "not a git repository",
		}
	}

	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, GitError{
			Operation: "branch",
			Message:   "failed to list branches: " + err.Error(),
		}
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
		return "", GitError{
			Operation: "branch",
			Message:   "not a git repository",
		}
	}

	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", GitError{
			Operation: "branch",
			Message:   "failed to get current branch: " + err.Error(),
		}
	}

	currentBranch := strings.TrimSpace(string(output))
	if currentBranch == "" {
		return "", GitError{
			Operation: "branch",
			Message:   "no current branch (detached HEAD?)",
		}
	}

	return currentBranch, nil
}

// CreateBranch creates a new Git branch from the specified base branch
func (g *LocalGitRepository) CreateBranch(name, baseBranch string) error {
	if !g.IsGitRepository() {
		return GitError{
			Operation: "branch",
			Message:   "not a git repository",
		}
	}

	if name == "" {
		return GitError{
			Operation: "branch",
			Message:   "branch name cannot be empty",
		}
	}

	if baseBranch == "" {
		return GitError{
			Operation: "branch",
			Message:   "base branch cannot be empty",
		}
	}

	cmd := exec.Command("git", "branch", name, baseBranch)
	if err := cmd.Run(); err != nil {
		return GitError{
			Operation: "branch",
			Message:   "failed to create branch '" + name + "' from '" + baseBranch + "': " + err.Error(),
		}
	}

	return nil
}

// CheckoutBranch switches to the specified Git branch
func (g *LocalGitRepository) CheckoutBranch(name string) error {
	if !g.IsGitRepository() {
		return GitError{
			Operation: "checkout",
			Message:   "not a git repository",
		}
	}

	if name == "" {
		return GitError{
			Operation: "checkout",
			Message:   "branch name cannot be empty",
		}
	}

	cmd := exec.Command("git", "checkout", name)
	if err := cmd.Run(); err != nil {
		return GitError{
			Operation: "checkout",
			Message:   "failed to checkout branch '" + name + "': " + err.Error(),
		}
	}

	return nil
}