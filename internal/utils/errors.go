package utils

import (
	"fmt"
)

type AnalysisError struct {
	Module   string
	Message string
	Cause    error
}

func (e *AnalysisError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Module, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Module, e.Message)
}

func (e *AnalysisError) Unwrap() error {
	return e.Cause
}

func NewAnalysisError(module, message string, cause error) *AnalysisError {
	return &AnalysisError{
		Module:   module,
		Message: message,
		Cause:    cause,
	}
}

type CleaningError struct {
	Module    string
	Message   string
	Cause     error
	ItemPath  string
}

func (e *CleaningError) Error() string {
	if e.ItemPath != "" {
		return fmt.Sprintf("%s: %s (item: %s)", e.Module, e.Message, e.ItemPath)
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Module, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Module, e.Message)
}

func (e *CleaningError) Unwrap() error {
	return e.Cause
}

func NewCleaningError(module, message string, cause error) *CleaningError {
	return &CleaningError{
		Module:  module,
		Message: message,
		Cause:   cause,
	}
}


