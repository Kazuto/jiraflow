package branch

// BranchInfo represents information needed to generate a branch name
type BranchInfo struct {
	Type       string
	TicketID   string
	Title      string
	BaseBranch string
	FullName   string
}

// Generator interface defines branch name generation operations
type Generator interface {
	GenerateName(info BranchInfo) string
	ValidateName(name string) error
}