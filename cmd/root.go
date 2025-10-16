package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jiraflow",
	Short: "Interactive Git branch creation tool for Jira workflows",
	Long: `JiraFlow is a CLI tool that helps developers create Git branches 
following Jira ticket naming conventions with an interactive TUI interface.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add flags and configuration here
}