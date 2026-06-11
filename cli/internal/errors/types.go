package errors

import "fmt"

// ErrorType categorises errors for consistent handling and user messaging.
type ErrorType string

const (
	ErrorTypeProject       ErrorType = "project"
	ErrorTypeDocker        ErrorType = "docker"
	ErrorTypeTemplate      ErrorType = "template"
	ErrorTypeRegistry      ErrorType = "registry"
	ErrorTypeFile          ErrorType = "file"
	ErrorTypeValidation    ErrorType = "validation"
	ErrorTypeConfiguration ErrorType = "configuration"
	ErrorTypeNetwork       ErrorType = "network"
	ErrorTypeAuth          ErrorType = "auth"
)

// AtempoError is a structured error with type, operation context, and optional detail bag.
type AtempoError struct {
	Type      ErrorType
	Operation string
	Message   string
	Cause     error
	Details   map[string]interface{}
}

func (e *AtempoError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Operation, e.Message)
}

func (e *AtempoError) Unwrap() error { return e.Cause }

func (e *AtempoError) Is(target error) bool {
	if t, ok := target.(*AtempoError); ok {
		return e.Type == t.Type
	}
	return false
}

// WithDetail adds a key/value pair to the error's detail bag (fluent).
func (e *AtempoError) WithDetail(key string, value interface{}) *AtempoError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

func (e *AtempoError) GetDetail(key string) (interface{}, bool) {
	if e.Details == nil {
		return nil, false
	}
	v, ok := e.Details[key]
	return v, ok
}

func NewProjectError(operation, message string, cause error) *AtempoError {
	return &AtempoError{Type: ErrorTypeProject, Operation: operation, Message: message, Cause: cause}
}

func NewDockerError(operation, message string, cause error) *AtempoError {
	return &AtempoError{Type: ErrorTypeDocker, Operation: operation, Message: message, Cause: cause}
}

func NewTemplateError(operation, message string, cause error) *AtempoError {
	return &AtempoError{Type: ErrorTypeTemplate, Operation: operation, Message: message, Cause: cause}
}

func NewRegistryError(operation, message string, cause error) *AtempoError {
	return &AtempoError{Type: ErrorTypeRegistry, Operation: operation, Message: message, Cause: cause}
}

func NewFileError(operation, message string, cause error) *AtempoError {
	return &AtempoError{Type: ErrorTypeFile, Operation: operation, Message: message, Cause: cause}
}

func NewValidationError(operation, message string, cause error) *AtempoError {
	return &AtempoError{Type: ErrorTypeValidation, Operation: operation, Message: message, Cause: cause}
}

func NewConfigurationError(operation, message string, cause error) *AtempoError {
	return &AtempoError{Type: ErrorTypeConfiguration, Operation: operation, Message: message, Cause: cause}
}

func NewNetworkError(operation, message string, cause error) *AtempoError {
	return &AtempoError{Type: ErrorTypeNetwork, Operation: operation, Message: message, Cause: cause}
}

func NewAuthError(operation, message string, cause error) *AtempoError {
	return &AtempoError{Type: ErrorTypeAuth, Operation: operation, Message: message, Cause: cause}
}
