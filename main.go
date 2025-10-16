package main

import (
	"os"

	"jiraflow/cmd"
	"jiraflow/internal/errors"
)

// Version information - set during build time
var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	// Set version information for the CLI
	cmd.SetVersionInfo(version, buildTime)
	
	errorHandler := errors.NewErrorHandler()
	
	if err := cmd.Execute(); err != nil {
		exitCode := errorHandler.HandleError(err)
		os.Exit(exitCode)
	}
}

