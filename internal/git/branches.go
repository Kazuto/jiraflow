package git

// BranchInfo represents information about a Git branch
type BranchInfo struct {
	Name      string
	IsCurrent bool
	IsRemote  bool
}