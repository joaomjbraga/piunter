package utils

import (
	"fmt"
	"strings"
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

func NewItemCleaningError(module, itemPath, message string, cause error) *CleaningError {
	return &CleaningError{
		Module:   module,
		Message:  message,
		Cause:    cause,
		ItemPath: itemPath,
	}
}

type ErrorHandler struct {
	errors []error
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{errors: make([]error, 0)}
}

func (h *ErrorHandler) Add(err error) {
	if err != nil {
		h.errors = append(h.errors, err)
	}
}

func (h *ErrorHandler) AddIf(err error, condition bool) {
	if err != nil && condition {
		h.errors = append(h.errors, err)
	}
}

func (h *ErrorHandler) Errors() []error {
	return h.errors
}

func (h *ErrorHandler) HasErrors() bool {
	return len(h.errors) > 0
}

func (h *ErrorHandler) Error() string {
	if !h.HasErrors() {
		return ""
	}
	
	var msgs []string
	for _, err := range h.errors {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

func (h *ErrorHandler) ToStrings() []string {
	var strs []string
	for _, err := range h.errors {
		strs = append(strs, err.Error())
	}
	return strs
}