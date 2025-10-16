package git

import (
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