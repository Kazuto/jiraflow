package main

import (
	"os"

	"jiraflow/cmd"
	"jiraflow/internal/errors"
)

func main() {
	errorHandler := errors.NewErrorHandler()
	
	if err := cmd.Execute(); err != nil {
		exitCode := errorHandler.HandleError(err)
		os.Exit(exitCode)
	}
}

