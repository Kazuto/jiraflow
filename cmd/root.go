package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"jiraflow/internal/branch"
	"jiraflow/internal/config"
	"jiraflow/internal/git"
	"jiraflow/internal/jira"
	"jiraflow/internal/tui"
)

var (
	// Global flags
	interactive bool
	dryRun      bool
	
	// Non-interactive mode flags
	branchType   string
	baseBranch   string
	ticketNumber string
	ticketTitle  string
	
	// Version information (placeholders for build-time injection)
	appVersion   = "dev"      //nolint:unused // Set by build process
	appBuildTime = "unknown"  //nolint:unused // Set by build process
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jiraflow [flags]",
	Short: "Interactive Git branch creation tool for Jira workflows",
	Long: `JiraFlow is a CLI tool that helps developers create Git branches 
following Jira ticket naming conventions with an interactive TUI interface.

The tool supports both interactive and non-interactive modes:

Interactive Mode (default):
  Launch the TUI interface to interactively select branch type, base branch,
  and enter ticket information with real-time branch name preview.

Non-Interactive Mode:
  Provide all required information via command-line flags to create branches
  without user interaction, perfect for automation and scripting.

Examples:
  # Launch interactive mode (default)
  jiraflow

  # Create branch in non-interactive mode
  jiraflow --type feature --base main --ticket PROJ-123 --title "Add user authentication"

  # Preview branch name without creating (dry-run)
  jiraflow --dry-run --type hotfix --base develop --ticket PROJ-456 --title "Fix login bug"

  # Non-interactive with minimal flags (title fetched from Jira if available)
  jiraflow --type feature --ticket PROJ-789`,
	RunE: runJiraFlow,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo sets the version information for the application
func SetVersionInfo(version, buildTime string) {
	appVersion = version
	appBuildTime = buildTime
	rootCmd.Version = fmt.Sprintf("%s (built %s)", version, buildTime)
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", true, "Run in interactive mode (default)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Preview branch name without creating the branch")
	
	// Non-interactive mode flags
	rootCmd.Flags().StringVarP(&branchType, "type", "t", "", "Branch type (feature, hotfix, refactor, support)")
	rootCmd.Flags().StringVarP(&baseBranch, "base", "b", "", "Base branch to create new branch from (defaults to current branch)")
	rootCmd.Flags().StringVar(&ticketNumber, "ticket", "", "Jira ticket number (e.g., PROJ-123)")
	rootCmd.Flags().StringVar(&ticketTitle, "title", "", "Ticket title (optional, will fetch from Jira if not provided)")
	
	// Mark flags as mutually exclusive with interactive mode
	rootCmd.MarkFlagsMutuallyExclusive("interactive", "type")
	rootCmd.MarkFlagsMutuallyExclusive("interactive", "base")
	rootCmd.MarkFlagsMutuallyExclusive("interactive", "ticket")
	rootCmd.MarkFlagsMutuallyExclusive("interactive", "title")
	
	// Add help command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "help",
		Short: "Show detailed help and examples",
		Long: `Show detailed help information and usage examples for JiraFlow.

JiraFlow supports two modes of operation:

1. Interactive Mode (default):
   Launch a terminal user interface (TUI) to interactively create branches.
   This mode provides a guided workflow with real-time branch name preview.

2. Non-Interactive Mode:
   Create branches directly from command-line arguments, perfect for automation.

Exit Codes:
  0 - Success or user cancelled
  1 - General error
  2 - Configuration error
  3 - Git repository error
  4 - Jira integration error

Configuration:
  JiraFlow automatically creates a configuration file at:
  ~/.config/jiraflow/jiraflow.yaml

  You can customize branch types, naming conventions, and other settings
  by editing this file.`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	})
}

// runJiraFlow is the main entry point for the CLI command
func runJiraFlow(cmd *cobra.Command, args []string) error {
	// Load configuration
	configManager := config.NewFileConfigManager()
	cfg, err := configManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize Git repository
	gitRepo := git.NewLocalGitRepository()
	if !gitRepo.IsGitRepository() {
		return fmt.Errorf("current directory is not a Git repository")
	}

	// Determine mode based on flags
	isNonInteractive := branchType != "" || baseBranch != "" || ticketNumber != ""
	
	if isNonInteractive {
		// Force interactive to false if any non-interactive flags are provided
		interactive = false
		return runNonInteractiveMode(cfg, gitRepo)
	}

	if interactive {
		return runInteractiveMode(cfg, gitRepo)
	}

	// This should not be reached, but provide fallback
	return runInteractiveMode(cfg, gitRepo)
}

// runInteractiveMode launches the TUI interface
func runInteractiveMode(cfg *config.Config, gitRepo git.GitRepository) error {
	if dryRun {
		fmt.Println("Note: Dry-run mode is not applicable in interactive mode.")
		fmt.Println("Branch preview will be shown in the confirmation step.")
		fmt.Println()
	}
	
	// Display keyboard shortcuts help before launching TUI
	fmt.Println("🚀 JiraFlow Interactive Mode")
	fmt.Println()
	fmt.Println("Keyboard Shortcuts:")
	fmt.Println("  ↑/↓ or j/k    Navigate lists")
	fmt.Println("  /             Search (in branch selection)")
	fmt.Println("  Enter         Select/Confirm")
	fmt.Println("  Esc           Go back")
	fmt.Println("  q or Ctrl+C   Quit")
	fmt.Println()
	fmt.Println("Press any key to continue...")
	
	// Wait for user input before launching TUI
	var input string
	_, _ = fmt.Scanln(&input)
	
	// Launch TUI and handle any errors
	if err := tui.RunTUI(cfg, gitRepo); err != nil {
		return fmt.Errorf("TUI application failed: %w", err)
	}
	
	return nil
}

// runNonInteractiveMode handles non-interactive branch creation
func runNonInteractiveMode(cfg *config.Config, gitRepo git.GitRepository) error {
	// Validate required flags for non-interactive mode
	if err := validateNonInteractiveFlags(cfg); err != nil {
		return err
	}

	// Get current branch as default base if not specified
	if baseBranch == "" {
		currentBranch, err := gitRepo.GetCurrentBranch()
		if err != nil {
			return fmt.Errorf("failed to get current branch (ensure you're in a Git repository): %w", err)
		}
		baseBranch = currentBranch
		fmt.Printf("Using current branch '%s' as base branch\n", baseBranch)
	} else {
		// Validate that the specified base branch exists
		branches, err := gitRepo.GetLocalBranches()
		if err != nil {
			return fmt.Errorf("failed to list local branches: %w", err)
		}
		
		branchExists := false
		for _, branch := range branches {
			if branch == baseBranch {
				branchExists = true
				break
			}
		}
		
		if !branchExists {
			return fmt.Errorf("base branch '%s' does not exist locally\nAvailable branches: %s", 
				baseBranch, strings.Join(branches, ", "))
		}
	}

	// Fetch ticket title from Jira if not provided and ticket number is given
	if ticketTitle == "" && ticketNumber != "" {
		jiraClient := jira.NewCLIClient()
		if title, err := jiraClient.GetTicketTitle(ticketNumber); err == nil {
			ticketTitle = title
			fmt.Printf("Fetched title from Jira: %s\n", ticketTitle)
		} else {
			fmt.Printf("Warning: Could not fetch title from Jira: %v\n", err)
			fmt.Println("Proceeding without title...")
		}
	}

	// Generate branch name
	sanitizer := branch.NewBranchSanitizer()
	generator := branch.NewBranchGenerator(sanitizer)
	branchInfo := branch.BranchInfo{
		Type:     branchType,
		TicketID: ticketNumber,
		Title:    ticketTitle,
	}
	branchName := generator.GenerateNameWithConfig(branchInfo, branch.GeneratorConfigFromAppConfig(
		cfg.MaxBranchLength,
		cfg.Sanitization.Separator,
		cfg.Sanitization.Lowercase,
		cfg.Sanitization.RemoveUmlauts,
	))

	// Display branch information
	fmt.Printf("\nBranch Information:\n")
	fmt.Printf("  Type: %s\n", branchType)
	fmt.Printf("  Base Branch: %s\n", baseBranch)
	fmt.Printf("  Ticket: %s\n", ticketNumber)
	if ticketTitle != "" {
		fmt.Printf("  Title: %s\n", ticketTitle)
	}
	fmt.Printf("  Generated Branch: %s\n", branchName)

	if dryRun {
		fmt.Printf("\n✓ Dry-run complete. Branch '%s' would be created from '%s'\n", branchName, baseBranch)
		return nil
	}

	// Check if branch already exists
	branches, err := gitRepo.GetLocalBranches()
	if err != nil {
		return fmt.Errorf("failed to check existing branches: %w", err)
	}
	
	for _, branch := range branches {
		if branch == branchName {
			return fmt.Errorf("branch '%s' already exists\nUse a different ticket number or title", branchName)
		}
	}

	// Create the branch
	fmt.Printf("\nCreating branch '%s' from '%s'...\n", branchName, baseBranch)
	if err := gitRepo.CreateBranch(branchName, baseBranch); err != nil {
		return fmt.Errorf("failed to create branch '%s': %w\nEnsure the base branch '%s' exists and you have proper Git permissions", 
			branchName, err, baseBranch)
	}

	fmt.Printf("✓ Successfully created and checked out branch '%s'\n", branchName)
	return nil
}

// validateNonInteractiveFlags validates the required flags for non-interactive mode
func validateNonInteractiveFlags(cfg *config.Config) error {
	var errors []string

	// Validate branch type
	if branchType == "" {
		errors = append(errors, "branch type is required (use --type flag)")
		
		// Provide helpful suggestion
		validTypes := make([]string, 0, len(cfg.BranchTypes))
		for t := range cfg.BranchTypes {
			validTypes = append(validTypes, t)
		}
		if len(validTypes) > 0 {
			errors = append(errors, fmt.Sprintf("  Available types: %s", strings.Join(validTypes, ", ")))
		}
	} else {
		// Check if branch type is valid
		validTypes := make([]string, 0, len(cfg.BranchTypes))
		for t := range cfg.BranchTypes {
			validTypes = append(validTypes, t)
		}
		
		isValid := false
		for _, validType := range validTypes {
			if branchType == validType {
				isValid = true
				break
			}
		}
		
		if !isValid {
			errors = append(errors, fmt.Sprintf("invalid branch type '%s'", branchType))
			errors = append(errors, fmt.Sprintf("  Valid types: %s", strings.Join(validTypes, ", ")))
		}
	}

	// Validate ticket number
	if ticketNumber == "" {
		errors = append(errors, "ticket number is required (use --ticket flag)")
		errors = append(errors, "  Example: --ticket PROJ-123")
	} else {
		// Basic ticket number format validation
		if len(ticketNumber) < 3 || !strings.Contains(ticketNumber, "-") {
			errors = append(errors, fmt.Sprintf("invalid ticket format '%s'", ticketNumber))
			errors = append(errors, "  Expected format: PROJECT-NUMBER (e.g., PROJ-123)")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed:\n  %s", strings.Join(errors, "\n  "))
	}

	return nil
}