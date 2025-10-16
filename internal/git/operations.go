package git

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