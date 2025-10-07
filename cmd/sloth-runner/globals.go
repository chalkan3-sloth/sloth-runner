package main

import (
	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	"github.com/chalkan3-sloth/sloth-runner/internal/taskrunner"
)

var surveyAsker taskrunner.SurveyAsker = &taskrunner.DefaultSurveyAsker{}

// GetSlothRunnerDataDir is a wrapper for backwards compatibility
func GetSlothRunnerDataDir() string {
	return config.GetDataDir()
}
