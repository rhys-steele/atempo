package errors

import (
	"errors"
	"fmt"
)

// IsType checks if an error is of a specific AtempoError type
func IsType(err error, errorType ErrorType) bool {
	var atempoErr *AtempoError
	if errors.As(err, &atempoErr) {
		return atempoErr.Type == errorType
	}
	return false
}

// GetType returns the ErrorType of an error, or empty string if not an AtempoError
func GetType(err error) ErrorType {
	var atempoErr *AtempoError
	if errors.As(err, &atempoErr) {
		return atempoErr.Type
	}
	return ""
}

// GetOperation returns the operation name of an error, or empty string if not an AtempoError
func GetOperation(err error) string {
	var atempoErr *AtempoError
	if errors.As(err, &atempoErr) {
		return atempoErr.Operation
	}
	return ""
}

// GetDetails returns the details map of an error, or nil if not an AtempoError
func GetDetails(err error) map[string]interface{} {
	var atempoErr *AtempoError
	if errors.As(err, &atempoErr) {
		return atempoErr.Details
	}
	return nil
}

// Wrap wraps an existing error as an AtempoError of the specified type
func Wrap(err error, errorType ErrorType, operation, message string) *AtempoError {
	return &AtempoError{
		Type:      errorType,
		Operation: operation,
		Message:   message,
		Cause:     err,
	}
}

// WrapWithDetails wraps an error and adds details
func WrapWithDetails(err error, errorType ErrorType, operation, message string, details map[string]interface{}) *AtempoError {
	return &AtempoError{
		Type:      errorType,
		Operation: operation,
		Message:   message,
		Cause:     err,
		Details:   details,
	}
}

// FormatUserMessage formats an error message for user display
// This provides a consistent, user-friendly error presentation
func FormatUserMessage(err error) string {
	var atempoErr *AtempoError
	if errors.As(err, &atempoErr) {
		switch atempoErr.Type {
		case ErrorTypeProject:
			return fmt.Sprintf("Project error: %s", atempoErr.Message)
		case ErrorTypeDocker:
			return fmt.Sprintf("Docker error: %s\nTip: Make sure Docker is running and accessible", atempoErr.Message)
		case ErrorTypeTemplate:
			return fmt.Sprintf("Template error: %s", atempoErr.Message)
		case ErrorTypeRegistry:
			return fmt.Sprintf("Registry error: %s", atempoErr.Message)
		case ErrorTypeFile:
			return fmt.Sprintf("File system error: %s", atempoErr.Message)
		case ErrorTypeValidation:
			return fmt.Sprintf("Validation error: %s", atempoErr.Message)
		case ErrorTypeConfiguration:
			return fmt.Sprintf("Configuration error: %s", atempoErr.Message)
		case ErrorTypeNetwork:
			return fmt.Sprintf("Network error: %s\nTip: Check your internet connection", atempoErr.Message)
		case ErrorTypeAuth:
			return fmt.Sprintf("Authentication error: %s", atempoErr.Message)
		default:
			return fmt.Sprintf("Error: %s", atempoErr.Message)
		}
	}
	return err.Error()
}

// GetRecoveryAction suggests a recovery action based on the error type
func GetRecoveryAction(err error) string {
	var atempoErr *AtempoError
	if errors.As(err, &atempoErr) {
		switch atempoErr.Type {
		case ErrorTypeDocker:
			return "Try running 'docker --version' to verify Docker is installed and accessible"
		case ErrorTypeProject:
			if projectName, ok := atempoErr.GetDetail("project_name"); ok {
				return fmt.Sprintf("Try 'atempo projects' to see available projects, or verify project '%v' exists", projectName)
			}
			return "Try 'atempo projects' to see available projects"
		case ErrorTypeTemplate:
			return "Try 'atempo create --help' to see available frameworks and versions"
		case ErrorTypeRegistry:
			return "Try 'atempo reset' to reinitialize the project registry"
		case ErrorTypeFile:
			if filePath, ok := atempoErr.GetDetail("file_path"); ok {
				return fmt.Sprintf("Check that file/directory exists and is accessible: %v", filePath)
			}
			return "Check file permissions and paths"
		case ErrorTypeValidation:
			return "Check the command syntax and arguments"
		case ErrorTypeConfiguration:
			return "Check the project's atempo.json configuration file"
		case ErrorTypeNetwork:
			return "Check your internet connection and proxy settings"
		case ErrorTypeAuth:
			return "Check your authentication credentials and permissions"
		default:
			return "Refer to documentation or run command with --help for more information"
		}
	}
	return "Run command with --help for more information"
}

// IsCritical determines if an error should stop execution vs allow graceful degradation
func IsCritical(err error) bool {
	var atempoErr *AtempoError
	if errors.As(err, &atempoErr) {
		// Most file and validation errors are critical
		switch atempoErr.Type {
		case ErrorTypeFile, ErrorTypeValidation, ErrorTypeAuth:
			return true
		case ErrorTypeDocker:
			// Docker errors might be recoverable depending on the operation
			return atempoErr.GetDetail("critical") == true
		case ErrorTypeNetwork:
			// Network errors are usually recoverable
			return false
		default:
			return true
		}
	}
	return true // Conservative approach - treat unknown errors as critical
}