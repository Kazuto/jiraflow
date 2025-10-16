package git

import (
	"os"
	"strings"
	"testing"
)

func TestFilterBranches(t *testing.T) {
	branches := []string{
		"main",
		"develop",
		"feature/user-auth",
		"feature/payment-system",
		"hotfix/critical-bug",
		"release/v1.2.0",
	}

	tests := []struct {
		name       string
		searchTerm string
		wantCount  int
		wantFirst  string
		hasResults bool
	}{
		{
			name:       "empty search returns all branches",
			searchTerm: "",
			wantCount:  6,
			wantFirst:  "main",
			hasResults: true,
		},
		{
			name:       "exact match",
			searchTerm: "main",
			wantCount:  1,
			wantFirst:  "main",
			hasResults: true,
		},
		{
			name:       "partial match",
			searchTerm: "feature",
			wantCount:  2,
			wantFirst:  "feature/user-auth",
			hasResults: true,
		},
		{
			name:       "fuzzy match",
			searchTerm: "ftr",
			wantCount:  2,
			wantFirst:  "feature/user-auth",
			hasResults: true,
		},
		{
			name:       "no matches",
			searchTerm: "nonexistent",
			wantCount:  0,
			wantFirst:  "",
			hasResults: false,
		},
		{
			name:       "case insensitive",
			searchTerm: "MAIN",
			wantCount:  1,
			wantFirst:  "main",
			hasResults: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterBranches(branches, tt.searchTerm)
			
			if len(result.Branches) != tt.wantCount {
				t.Errorf("FilterBranches() got %d branches, want %d", len(result.Branches), tt.wantCount)
			}
			
			if result.HasResults != tt.hasResults {
				t.Errorf("FilterBranches() HasResults = %v, want %v", result.HasResults, tt.hasResults)
			}
			
			if tt.wantCount > 0 && result.Branches[0] != tt.wantFirst {
				t.Errorf("FilterBranches() first result = %v, want %v", result.Branches[0], tt.wantFirst)
			}
			
			if result.SearchTerm != tt.searchTerm {
				t.Errorf("FilterBranches() SearchTerm = %v, want %v", result.SearchTerm, tt.searchTerm)
			}
		})
	}
}

func TestFilterBranchesRealtime(t *testing.T) {
	branches := []string{
		"main",
		"develop",
		"feature/user-authentication",
		"feature/payment-gateway",
		"hotfix/auth-bug",
	}

	tests := []struct {
		name       string
		searchTerm string
		wantCount  int
		hasResults bool
	}{
		{
			name:       "exact substring match prioritized",
			searchTerm: "auth",
			wantCount:  2,
			hasResults: true,
		},
		{
			name:       "fuzzy match works",
			searchTerm: "ftr",
			wantCount:  2,
			hasResults: true,
		},
		{
			name:       "no matches returns empty",
			searchTerm: "xyz",
			wantCount:  0,
			hasResults: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterBranchesRealtime(branches, tt.searchTerm)
			
			if len(result.Branches) != tt.wantCount {
				t.Errorf("FilterBranchesRealtime() got %d branches, want %d", len(result.Branches), tt.wantCount)
			}
			
			if result.HasResults != tt.hasResults {
				t.Errorf("FilterBranchesRealtime() HasResults = %v, want %v", result.HasResults, tt.hasResults)
			}
		})
	}
}

func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		searchTerm string
		want       bool
	}{
		{
			name:       "exact match",
			branchName: "main",
			searchTerm: "main",
			want:       true,
		},
		{
			name:       "fuzzy match in order",
			branchName: "feature/user-auth",
			searchTerm: "ftr",
			want:       true,
		},
		{
			name:       "fuzzy match out of order",
			branchName: "feature/user-auth",
			searchTerm: "rft",
			want:       false,
		},
		{
			name:       "partial match",
			branchName: "develop",
			searchTerm: "dev",
			want:       true,
		},
		{
			name:       "no match",
			branchName: "main",
			searchTerm: "xyz",
			want:       false,
		},
		{
			name:       "empty search term",
			branchName: "main",
			searchTerm: "",
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fuzzyMatch(tt.branchName, tt.searchTerm)
			if got != tt.want {
				t.Errorf("fuzzyMatch(%q, %q) = %v, want %v", tt.branchName, tt.searchTerm, got, tt.want)
			}
		})
	}
}

func TestBranchSearchResult_GetEmptySearchMessage(t *testing.T) {
	tests := []struct {
		name       string
		searchTerm string
		want       string
	}{
		{
			name:       "empty search term",
			searchTerm: "",
			want:       "No branches available",
		},
		{
			name:       "with search term",
			searchTerm: "nonexistent",
			want:       "No branches found matching 'nonexistent'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BranchSearchResult{SearchTerm: tt.searchTerm}
			got := result.GetEmptySearchMessage()
			if got != tt.want {
				t.Errorf("GetEmptySearchMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBranchSearchResult_GetSearchSummary(t *testing.T) {
	tests := []struct {
		name       string
		searchTerm string
		branches   []string
		want       string
	}{
		{
			name:       "no search term",
			searchTerm: "",
			branches:   []string{"main", "develop"},
			want:       "",
		},
		{
			name:       "no results",
			searchTerm: "xyz",
			branches:   []string{},
			want:       "No branches found matching 'xyz'",
		},
		{
			name:       "one result",
			searchTerm: "main",
			branches:   []string{"main"},
			want:       "1 branch found",
		},
		{
			name:       "multiple results",
			searchTerm: "feature",
			branches:   []string{"feature/auth", "feature/payment"},
			want:       "2 branches found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BranchSearchResult{
				SearchTerm: tt.searchTerm,
				Branches:   tt.branches,
			}
			got := result.GetSearchSummary()
			if got != tt.want {
				t.Errorf("GetSearchSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}

// MockGitRepository implements GitRepository for testing
type MockGitRepository struct {
	branches      []string
	currentBranch string
	isGitRepo     bool
	shouldError   bool
	errorMessage  string
}

func NewMockGitRepository() *MockGitRepository {
	return &MockGitRepository{
		branches: []string{
			"main",
			"develop", 
			"feature/user-authentication",
			"feature/payment-gateway",
			"hotfix/critical-bug",
			"release/v1.2.0",
		},
		currentBranch: "main",
		isGitRepo:     true,
		shouldError:   false,
	}
}

func (m *MockGitRepository) SetBranches(branches []string) {
	m.branches = branches
}

func (m *MockGitRepository) SetCurrentBranch(branch string) {
	m.currentBranch = branch
}

func (m *MockGitRepository) SetIsGitRepo(isRepo bool) {
	m.isGitRepo = isRepo
}

func (m *MockGitRepository) SetShouldError(shouldError bool, message string) {
	m.shouldError = shouldError
	m.errorMessage = message
}

func (m *MockGitRepository) IsGitRepository() bool {
	return m.isGitRepo
}

func (m *MockGitRepository) GetLocalBranches() ([]string, error) {
	if !m.isGitRepo {
		return nil, GitError{
			Operation: "branch",
			Message:   "not a git repository",
		}
	}
	
	if m.shouldError {
		return nil, GitError{
			Operation: "branch",
			Message:   m.errorMessage,
		}
	}
	
	return m.branches, nil
}

func (m *MockGitRepository) GetCurrentBranch() (string, error) {
	if !m.isGitRepo {
		return "", GitError{
			Operation: "branch",
			Message:   "not a git repository",
		}
	}
	
	if m.shouldError {
		return "", GitError{
			Operation: "branch",
			Message:   m.errorMessage,
		}
	}
	
	return m.currentBranch, nil
}

func (m *MockGitRepository) CreateBranch(name, baseBranch string) error {
	if !m.isGitRepo {
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
	
	if m.shouldError {
		return GitError{
			Operation: "branch",
			Message:   m.errorMessage,
		}
	}
	
	// Add the new branch to our mock list
	m.branches = append(m.branches, name)
	return nil
}

func (m *MockGitRepository) CheckoutBranch(name string) error {
	if !m.isGitRepo {
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
	
	if m.shouldError {
		return GitError{
			Operation: "checkout",
			Message:   m.errorMessage,
		}
	}
	
	// Check if branch exists
	branchExists := false
	for _, branch := range m.branches {
		if branch == name {
			branchExists = true
			break
		}
	}
	
	if !branchExists {
		return GitError{
			Operation: "checkout",
			Message:   "branch '" + name + "' does not exist",
		}
	}
	
	m.currentBranch = name
	return nil
}

func (m *MockGitRepository) SearchBranches(searchTerm string) (BranchSearchResult, error) {
	branches, err := m.GetLocalBranches()
	if err != nil {
		return BranchSearchResult{}, err
	}
	
	result := FilterBranchesRealtime(branches, searchTerm)
	return result, nil
}

// Test Git operations with mock repository
func TestLocalGitRepository_GetLocalBranches(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockGitRepository)
		wantCount   int
		wantError   bool
		errorSubstr string
	}{
		{
			name: "successful branch listing",
			setupMock: func(m *MockGitRepository) {
				m.SetBranches([]string{"main", "develop", "feature/test"})
			},
			wantCount: 3,
			wantError: false,
		},
		{
			name: "not a git repository",
			setupMock: func(m *MockGitRepository) {
				m.SetIsGitRepo(false)
			},
			wantCount:   0,
			wantError:   true,
			errorSubstr: "not a git repository",
		},
		{
			name: "git command fails",
			setupMock: func(m *MockGitRepository) {
				m.SetShouldError(true, "failed to list branches")
			},
			wantCount:   0,
			wantError:   true,
			errorSubstr: "failed to list branches",
		},
		{
			name: "empty repository",
			setupMock: func(m *MockGitRepository) {
				m.SetBranches([]string{})
			},
			wantCount: 0,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockGitRepository()
			tt.setupMock(mock)
			
			branches, err := mock.GetLocalBranches()
			
			if tt.wantError {
				if err == nil {
					t.Errorf("GetLocalBranches() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("GetLocalBranches() error = %v, want error containing %v", err, tt.errorSubstr)
				}
			} else {
				if err != nil {
					t.Errorf("GetLocalBranches() unexpected error = %v", err)
				}
				if len(branches) != tt.wantCount {
					t.Errorf("GetLocalBranches() got %d branches, want %d", len(branches), tt.wantCount)
				}
			}
		})
	}
}

func TestLocalGitRepository_GetCurrentBranch(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockGitRepository)
		wantBranch  string
		wantError   bool
		errorSubstr string
	}{
		{
			name: "successful current branch detection",
			setupMock: func(m *MockGitRepository) {
				m.SetCurrentBranch("feature/test")
			},
			wantBranch: "feature/test",
			wantError:  false,
		},
		{
			name: "not a git repository",
			setupMock: func(m *MockGitRepository) {
				m.SetIsGitRepo(false)
			},
			wantBranch:  "",
			wantError:   true,
			errorSubstr: "not a git repository",
		},
		{
			name: "git command fails",
			setupMock: func(m *MockGitRepository) {
				m.SetShouldError(true, "failed to get current branch")
			},
			wantBranch:  "",
			wantError:   true,
			errorSubstr: "failed to get current branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockGitRepository()
			tt.setupMock(mock)
			
			branch, err := mock.GetCurrentBranch()
			
			if tt.wantError {
				if err == nil {
					t.Errorf("GetCurrentBranch() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("GetCurrentBranch() error = %v, want error containing %v", err, tt.errorSubstr)
				}
			} else {
				if err != nil {
					t.Errorf("GetCurrentBranch() unexpected error = %v", err)
				}
				if branch != tt.wantBranch {
					t.Errorf("GetCurrentBranch() = %v, want %v", branch, tt.wantBranch)
				}
			}
		})
	}
}

func TestLocalGitRepository_CreateBranch(t *testing.T) {
	tests := []struct {
		name        string
		branchName  string
		baseBranch  string
		setupMock   func(*MockGitRepository)
		wantError   bool
		errorSubstr string
	}{
		{
			name:       "successful branch creation",
			branchName: "feature/new-feature",
			baseBranch: "main",
			setupMock:  func(m *MockGitRepository) {},
			wantError:  false,
		},
		{
			name:        "empty branch name",
			branchName:  "",
			baseBranch:  "main",
			setupMock:   func(m *MockGitRepository) {},
			wantError:   true,
			errorSubstr: "branch name cannot be empty",
		},
		{
			name:        "empty base branch",
			branchName:  "feature/test",
			baseBranch:  "",
			setupMock:   func(m *MockGitRepository) {},
			wantError:   true,
			errorSubstr: "base branch cannot be empty",
		},
		{
			name:       "not a git repository",
			branchName: "feature/test",
			baseBranch: "main",
			setupMock: func(m *MockGitRepository) {
				m.SetIsGitRepo(false)
			},
			wantError:   true,
			errorSubstr: "not a git repository",
		},
		{
			name:       "git command fails",
			branchName: "feature/test",
			baseBranch: "main",
			setupMock: func(m *MockGitRepository) {
				m.SetShouldError(true, "failed to create branch")
			},
			wantError:   true,
			errorSubstr: "failed to create branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockGitRepository()
			tt.setupMock(mock)
			
			err := mock.CreateBranch(tt.branchName, tt.baseBranch)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("CreateBranch() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("CreateBranch() error = %v, want error containing %v", err, tt.errorSubstr)
				}
			} else {
				if err != nil {
					t.Errorf("CreateBranch() unexpected error = %v", err)
				}
				// Verify branch was added to mock
				branches, _ := mock.GetLocalBranches()
				found := false
				for _, branch := range branches {
					if branch == tt.branchName {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("CreateBranch() branch %v was not added to repository", tt.branchName)
				}
			}
		})
	}
}

func TestLocalGitRepository_CheckoutBranch(t *testing.T) {
	tests := []struct {
		name        string
		branchName  string
		setupMock   func(*MockGitRepository)
		wantError   bool
		errorSubstr string
	}{
		{
			name:       "successful branch checkout",
			branchName: "develop",
			setupMock: func(m *MockGitRepository) {
				m.SetBranches([]string{"main", "develop", "feature/test"})
			},
			wantError: false,
		},
		{
			name:        "empty branch name",
			branchName:  "",
			setupMock:   func(m *MockGitRepository) {},
			wantError:   true,
			errorSubstr: "branch name cannot be empty",
		},
		{
			name:       "branch does not exist",
			branchName: "nonexistent",
			setupMock: func(m *MockGitRepository) {
				m.SetBranches([]string{"main", "develop"})
			},
			wantError:   true,
			errorSubstr: "branch 'nonexistent' does not exist",
		},
		{
			name:       "not a git repository",
			branchName: "main",
			setupMock: func(m *MockGitRepository) {
				m.SetIsGitRepo(false)
			},
			wantError:   true,
			errorSubstr: "not a git repository",
		},
		{
			name:       "git command fails",
			branchName: "main",
			setupMock: func(m *MockGitRepository) {
				m.SetShouldError(true, "failed to checkout")
			},
			wantError:   true,
			errorSubstr: "failed to checkout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockGitRepository()
			tt.setupMock(mock)
			
			err := mock.CheckoutBranch(tt.branchName)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("CheckoutBranch() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("CheckoutBranch() error = %v, want error containing %v", err, tt.errorSubstr)
				}
			} else {
				if err != nil {
					t.Errorf("CheckoutBranch() unexpected error = %v", err)
				}
				// Verify current branch was updated
				currentBranch, _ := mock.GetCurrentBranch()
				if currentBranch != tt.branchName {
					t.Errorf("CheckoutBranch() current branch = %v, want %v", currentBranch, tt.branchName)
				}
			}
		})
	}
}

func TestLocalGitRepository_SearchBranches(t *testing.T) {
	tests := []struct {
		name        string
		searchTerm  string
		setupMock   func(*MockGitRepository)
		wantCount   int
		wantError   bool
		errorSubstr string
	}{
		{
			name:       "successful search with results",
			searchTerm: "feature",
			setupMock: func(m *MockGitRepository) {
				m.SetBranches([]string{"main", "feature/auth", "feature/payment", "hotfix/bug"})
			},
			wantCount: 2,
			wantError: false,
		},
		{
			name:       "search with no results",
			searchTerm: "nonexistent",
			setupMock: func(m *MockGitRepository) {
				m.SetBranches([]string{"main", "develop"})
			},
			wantCount: 0,
			wantError: false,
		},
		{
			name:       "empty search returns all branches",
			searchTerm: "",
			setupMock: func(m *MockGitRepository) {
				m.SetBranches([]string{"main", "develop", "feature/test"})
			},
			wantCount: 3,
			wantError: false,
		},
		{
			name:       "search fails when not git repository",
			searchTerm: "main",
			setupMock: func(m *MockGitRepository) {
				m.SetIsGitRepo(false)
			},
			wantCount:   0,
			wantError:   true,
			errorSubstr: "not a git repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockGitRepository()
			tt.setupMock(mock)
			
			result, err := mock.SearchBranches(tt.searchTerm)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("SearchBranches() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("SearchBranches() error = %v, want error containing %v", err, tt.errorSubstr)
				}
			} else {
				if err != nil {
					t.Errorf("SearchBranches() unexpected error = %v", err)
				}
				if len(result.Branches) != tt.wantCount {
					t.Errorf("SearchBranches() got %d branches, want %d", len(result.Branches), tt.wantCount)
				}
				if result.SearchTerm != tt.searchTerm {
					t.Errorf("SearchBranches() SearchTerm = %v, want %v", result.SearchTerm, tt.searchTerm)
				}
			}
		})
	}
}

// Test error handling for non-Git directories with real LocalGitRepository
func TestLocalGitRepository_NonGitDirectory(t *testing.T) {
	// Create a temporary directory that is not a Git repository
	tempDir := t.TempDir()
	
	// Change to the temporary directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore original directory: %v", err)
		}
	}()
	
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	
	repo := NewLocalGitRepository()
	
	// Test IsGitRepository
	if repo.IsGitRepository() {
		t.Errorf("IsGitRepository() = true, want false for non-git directory")
	}
	
	// Test GetLocalBranches
	branches, err := repo.GetLocalBranches()
	if err == nil {
		t.Errorf("GetLocalBranches() expected error but got none")
	}
	if branches != nil {
		t.Errorf("GetLocalBranches() = %v, want nil for non-git directory", branches)
	}
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Errorf("GetLocalBranches() error = %v, want error containing 'not a git repository'", err)
	}
	
	// Test GetCurrentBranch
	currentBranch, err := repo.GetCurrentBranch()
	if err == nil {
		t.Errorf("GetCurrentBranch() expected error but got none")
	}
	if currentBranch != "" {
		t.Errorf("GetCurrentBranch() = %v, want empty string for non-git directory", currentBranch)
	}
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Errorf("GetCurrentBranch() error = %v, want error containing 'not a git repository'", err)
	}
	
	// Test CreateBranch
	err = repo.CreateBranch("test-branch", "main")
	if err == nil {
		t.Errorf("CreateBranch() expected error but got none")
	}
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Errorf("CreateBranch() error = %v, want error containing 'not a git repository'", err)
	}
	
	// Test CheckoutBranch
	err = repo.CheckoutBranch("main")
	if err == nil {
		t.Errorf("CheckoutBranch() expected error but got none")
	}
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Errorf("CheckoutBranch() error = %v, want error containing 'not a git repository'", err)
	}
}

// Test with a real Git repository (if available)
func TestLocalGitRepository_RealGitRepo(t *testing.T) {
	// Skip this test if we're not in a Git repository
	repo := NewLocalGitRepository()
	if !repo.IsGitRepository() {
		t.Skip("Skipping test: not in a Git repository")
	}
	
	// Test GetLocalBranches
	branches, err := repo.GetLocalBranches()
	if err != nil {
		t.Errorf("GetLocalBranches() unexpected error = %v", err)
	}
	if branches == nil {
		t.Errorf("GetLocalBranches() = nil, want non-nil slice")
	}
	
	// Test GetCurrentBranch
	currentBranch, err := repo.GetCurrentBranch()
	if err != nil {
		// In CI environments (especially during releases), we might be in detached HEAD state
		// This is acceptable, so we'll just log it and continue
		t.Logf("GetCurrentBranch() returned error (likely detached HEAD): %v", err)
	} else if currentBranch != "" {
		// If we have a current branch, verify it's in the list of branches
		found := false
		for _, branch := range branches {
			if branch == currentBranch {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Current branch %v not found in branch list %v", currentBranch, branches)
		}
	} else {
		// Empty current branch might indicate detached HEAD, which is acceptable in CI
		t.Logf("GetCurrentBranch() returned empty string (likely detached HEAD)")
	}
}

// Test branch name patterns and edge cases
func TestBranchNamePatterns(t *testing.T) {
	branches := []string{
		"main",
		"master", 
		"develop",
		"dev",
		"feature/JIRA-123-user-authentication",
		"feature/payment-gateway-integration",
		"hotfix/critical-security-patch",
		"release/v1.2.0",
		"release/v2.0.0-beta",
		"bugfix/fix-login-issue",
		"chore/update-dependencies",
		"docs/update-readme",
		"test/add-integration-tests",
	}
	
	tests := []struct {
		name       string
		searchTerm string
		wantCount  int
		wantFirst  string
	}{
		{
			name:       "search by prefix",
			searchTerm: "feature",
			wantCount:  2,
			wantFirst:  "feature/JIRA-123-user-authentication",
		},
		{
			name:       "search by ticket number",
			searchTerm: "JIRA-123",
			wantCount:  1,
			wantFirst:  "feature/JIRA-123-user-authentication",
		},
		{
			name:       "search by version",
			searchTerm: "v1.2",
			wantCount:  1,
			wantFirst:  "release/v1.2.0",
		},
		{
			name:       "fuzzy search",
			searchTerm: "ftr",
			wantCount:  3, // feature branches + hotfix (fuzzy match)
			wantFirst:  "feature/JIRA-123-user-authentication",
		},
		{
			name:       "case insensitive search",
			searchTerm: "MAIN",
			wantCount:  2, // main + feature/payment-gateway-integration (fuzzy match)
			wantFirst:  "main",
		},
		{
			name:       "search with hyphens",
			searchTerm: "user-auth",
			wantCount:  1,
			wantFirst:  "feature/JIRA-123-user-authentication",
		},
		{
			name:       "search multiple word types",
			searchTerm: "update",
			wantCount:  2, // chore/update-dependencies and docs/update-readme
			wantFirst:  "chore/update-dependencies",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterBranchesRealtime(branches, tt.searchTerm)
			
			if len(result.Branches) != tt.wantCount {
				t.Errorf("FilterBranchesRealtime() got %d branches, want %d. Branches: %v", 
					len(result.Branches), tt.wantCount, result.Branches)
			}
			
			if tt.wantCount > 0 && result.Branches[0] != tt.wantFirst {
				t.Errorf("FilterBranchesRealtime() first result = %v, want %v", 
					result.Branches[0], tt.wantFirst)
			}
		})
	}
}

// Test GitError type
func TestGitError(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		message   string
		wantError string
	}{
		{
			name:      "branch operation error",
			operation: "branch",
			message:   "not a git repository",
			wantError: "git branch: not a git repository",
		},
		{
			name:      "checkout operation error",
			operation: "checkout",
			message:   "branch does not exist",
			wantError: "git checkout: branch does not exist",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GitError{
				Operation: tt.operation,
				Message:   tt.message,
			}
			
			if err.Error() != tt.wantError {
				t.Errorf("GitError.Error() = %v, want %v", err.Error(), tt.wantError)
			}
		})
	}
}