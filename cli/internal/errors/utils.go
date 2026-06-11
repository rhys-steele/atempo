package errors

import "errors"

// IsType reports whether err is an AtempoError of the given type.
func IsType(err error, t ErrorType) bool {
	var e *AtempoError
	return errors.As(err, &e) && e.Type == t
}

// GetDetails returns the details map, or nil if err is not an AtempoError.
func GetDetails(err error) map[string]interface{} {
	var e *AtempoError
	if errors.As(err, &e) {
		return e.Details
	}
	return nil
}

// Wrap wraps err as an AtempoError.
func Wrap(err error, t ErrorType, operation, message string) *AtempoError {
	return &AtempoError{Type: t, Operation: operation, Message: message, Cause: err}
}

// FormatUserMessage returns a user-facing description of err.
func FormatUserMessage(err error) string {
	var e *AtempoError
	if !errors.As(err, &e) {
		return err.Error()
	}
	switch e.Type {
	case ErrorTypeDocker:
		return "Docker error: " + e.Message + "\nTip: make sure Docker is running"
	case ErrorTypeNetwork:
		return "Network error: " + e.Message + "\nTip: check your internet connection"
	default:
		if e.Type != "" {
			return string(e.Type) + " error: " + e.Message
		}
		return "error: " + e.Message
	}
}

// GetRecoveryAction returns a suggested recovery step for err.
func GetRecoveryAction(err error) string {
	var e *AtempoError
	if !errors.As(err, &e) {
		return ""
	}
	switch e.Type {
	case ErrorTypeDocker:
		return "run 'docker --version' to verify Docker is installed and accessible"
	case ErrorTypeProject:
		if name, ok := e.GetDetail("project_name"); ok {
			if s, ok := name.(string); ok {
				return "run 'atempo projects' to see available projects, or check that '" + s + "' exists"
			}
		}
		return "run 'atempo projects' to see available projects"
	case ErrorTypeRegistry:
		return "run 'atempo reset' to reinitialise the project registry"
	case ErrorTypeValidation:
		return "check the command syntax — run with --help for usage"
	case ErrorTypeConfiguration:
		return "check the project's atempo.json configuration file"
	case ErrorTypeNetwork:
		return "check your internet connection and proxy settings"
	case ErrorTypeAuth:
		return "check your authentication credentials"
	default:
		return ""
	}
}

// IsCritical reports whether err should halt execution.
func IsCritical(err error) bool {
	var e *AtempoError
	if !errors.As(err, &e) {
		return true
	}
	switch e.Type {
	case ErrorTypeNetwork, ErrorTypeDocker:
		return false
	default:
		return true
	}
}
