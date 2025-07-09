package errors

import (
	"fmt"
)

// ErrorType represents the type of error for consistent categorization
type ErrorType string

const (
	// ErrorTypeProject represents project-related errors
	ErrorTypeProject ErrorType = "project"
	
	// ErrorTypeDocker represents Docker-related errors
	ErrorTypeDocker ErrorType = "docker"
	
	// ErrorTypeTemplate represents template-related errors
	ErrorTypeTemplate ErrorType = "template"
	
	// ErrorTypeRegistry represents registry-related errors
	ErrorTypeRegistry ErrorType = "registry"
	
	// ErrorTypeFile represents file system-related errors
	ErrorTypeFile ErrorType = "file"
	
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "validation"
	
	// ErrorTypeConfiguration represents configuration errors
	ErrorTypeConfiguration ErrorType = "configuration"
	
	// ErrorTypeNetwork represents network-related errors
	ErrorTypeNetwork ErrorType = "network"
	
	// ErrorTypeAuth represents authentication/authorization errors
	ErrorTypeAuth ErrorType = "auth"
)

// AtempoError represents a structured error with context and type information
type AtempoError struct {
	Type      ErrorType
	Operation string
	Message   string
	Cause     error
	Details   map[string]interface{}
}

// Error implements the error interface
func (e *AtempoError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s operation failed: %s: %v", e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s operation failed: %s", e.Operation, e.Message)
}

// Unwrap returns the underlying error
func (e *AtempoError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches the target type
func (e *AtempoError) Is(target error) bool {
	if targetErr, ok := target.(*AtempoError); ok {
		return e.Type == targetErr.Type
	}
	return false
}

// WithDetail adds a detail key-value pair to the error
func (e *AtempoError) WithDetail(key string, value interface{}) *AtempoError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// GetDetail retrieves a detail value by key
func (e *AtempoError) GetDetail(key string) (interface{}, bool) {
	if e.Details == nil {
		return nil, false
	}
	value, exists := e.Details[key]
	return value, exists
}

// Common error creation functions

// NewProjectError creates a project-related error
func NewProjectError(operation, message string, cause error) *AtempoError {
	return &AtempoError{
		Type:      ErrorTypeProject,
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}

// NewDockerError creates a Docker-related error
func NewDockerError(operation, message string, cause error) *AtempoError {
	return &AtempoError{
		Type:      ErrorTypeDocker,
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}

// NewTemplateError creates a template-related error
func NewTemplateError(operation, message string, cause error) *AtempoError {
	return &AtempoError{
		Type:      ErrorTypeTemplate,
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}

// NewRegistryError creates a registry-related error
func NewRegistryError(operation, message string, cause error) *AtempoError {
	return &AtempoError{
		Type:      ErrorTypeRegistry,
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}

// NewFileError creates a file system-related error
func NewFileError(operation, message string, cause error) *AtempoError {
	return &AtempoError{
		Type:      ErrorTypeFile,
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}

// NewValidationError creates a validation error
func NewValidationError(operation, message string, cause error) *AtempoError {
	return &AtempoError{
		Type:      ErrorTypeValidation,
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}

// NewConfigurationError creates a configuration error
func NewConfigurationError(operation, message string, cause error) *AtempoError {
	return &AtempoError{
		Type:      ErrorTypeConfiguration,
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}

// NewNetworkError creates a network-related error
func NewNetworkError(operation, message string, cause error) *AtempoError {
	return &AtempoError{
		Type:      ErrorTypeNetwork,
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}

// NewAuthError creates an authentication/authorization error
func NewAuthError(operation, message string, cause error) *AtempoError {
	return &AtempoError{
		Type:      ErrorTypeAuth,
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}