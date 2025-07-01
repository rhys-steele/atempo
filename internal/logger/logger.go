package logger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Logger provides project-specific logging with progress indicators
type Logger struct {
	ProjectName string
	LogFile     *os.File
	LogPath     string
	StartTime   time.Time
}

// StepStatus represents the status of a setup step
type StepStatus int

const (
	StepPending StepStatus = iota
	StepRunning
	StepComplete
	StepWarning
	StepError
)

// Step represents a setup step with timing and status
type Step struct {
	Name      string
	Status    StepStatus
	StartTime time.Time
	Duration  time.Duration
	Error     error
}

// New creates a new logger for a project
func New(projectName string) (*Logger, error) {
	// Create logs directory in the project or home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	logsDir := filepath.Join(homeDir, ".steele", "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create timestamped log file
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("%s_%s.log", projectName, timestamp)
	logPath := filepath.Join(logsDir, logFileName)

	logFile, err := os.Create(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	logger := &Logger{
		ProjectName: projectName,
		LogFile:     logFile,
		LogPath:     logPath,
		StartTime:   time.Now(),
	}

	// Write header to log file
	logger.writeHeader()

	return logger, nil
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.LogFile != nil {
		l.writeFooter()
		return l.LogFile.Close()
	}
	return nil
}

// writeHeader writes the log file header
func (l *Logger) writeHeader() {
	header := fmt.Sprintf(`
========================================
Steele Project Setup Log
========================================
Project: %s
Started: %s
========================================

`, l.ProjectName, l.StartTime.Format("2006-01-02 15:04:05"))
	l.LogFile.WriteString(header)
}

// writeFooter writes the log file footer
func (l *Logger) writeFooter() {
	duration := time.Since(l.StartTime)
	footer := fmt.Sprintf(`
========================================
Setup completed in %s
Log file: %s
========================================
`, duration.Round(time.Second), l.LogPath)
	l.LogFile.WriteString(footer)
}

// StartStep begins a new step and returns a Step instance
func (l *Logger) StartStep(name string) *Step {
	step := &Step{
		Name:      name,
		Status:    StepRunning,
		StartTime: time.Now(),
	}
	
	// Write to log file
	l.logf("[%s] STARTED: %s", step.StartTime.Format("15:04:05"), name)
	
	// Show progress indicator
	l.showProgress(step)
	
	return step
}

// CompleteStep marks a step as complete
func (l *Logger) CompleteStep(step *Step) {
	step.Status = StepComplete
	step.Duration = time.Since(step.StartTime)
	
	// Write to log file
	l.logf("[%s] COMPLETED: %s (took %s)", 
		time.Now().Format("15:04:05"), 
		step.Name, 
		step.Duration.Round(time.Millisecond))
	
	// Update progress indicator
	l.showProgress(step)
}

// WarningStep marks a step as completed with warnings
func (l *Logger) WarningStep(step *Step, warning string) {
	step.Status = StepWarning
	step.Duration = time.Since(step.StartTime)
	step.Error = fmt.Errorf("warning: %s", warning)
	
	// Write to log file
	l.logf("[%s] WARNING: %s (took %s) - %s", 
		time.Now().Format("15:04:05"), 
		step.Name, 
		step.Duration.Round(time.Millisecond),
		warning)
	
	// Update progress indicator
	l.showProgress(step)
}

// ErrorStep marks a step as failed
func (l *Logger) ErrorStep(step *Step, err error) {
	step.Status = StepError
	step.Duration = time.Since(step.StartTime)
	step.Error = err
	
	// Write to log file
	l.logf("[%s] ERROR: %s (took %s) - %s", 
		time.Now().Format("15:04:05"), 
		step.Name, 
		step.Duration.Round(time.Millisecond),
		err.Error())
	
	// Update progress indicator
	l.showProgress(step)
}

// RunCommand executes a command and captures its output to the log file
func (l *Logger) RunCommand(step *Step, cmd *exec.Cmd) error {
	// Log the command being executed
	l.logf("EXECUTING: %s", strings.Join(cmd.Args, " "))
	if cmd.Dir != "" {
		l.logf("WORKING DIR: %s", cmd.Dir)
	}
	
	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	
	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}
	
	// Capture output in goroutines
	done := make(chan struct{})
	go l.captureOutput(stdoutPipe, "STDOUT", done)
	go l.captureOutput(stderrPipe, "STDERR", done)
	
	// Wait for command completion
	err = cmd.Wait()
	
	// Wait for output capture to complete
	<-done
	<-done
	
	// Log completion
	if err != nil {
		l.logf("COMMAND FAILED: %s", err.Error())
		return err
	} else {
		l.logf("COMMAND COMPLETED SUCCESSFULLY")
	}
	
	return nil
}

// captureOutput captures command output and writes it to the log file
func (l *Logger) captureOutput(reader io.Reader, prefix string, done chan struct{}) {
	defer func() { done <- struct{}{} }()
	
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		l.logf("%s: %s", prefix, line)
	}
}

// logf writes a formatted message to the log file with timestamp
func (l *Logger) logf(format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05.000")
	message := fmt.Sprintf(format, args...)
	l.LogFile.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, message))
	l.LogFile.Sync() // Ensure it's written to disk immediately
}

// showProgress displays a progress indicator for the current step
func (l *Logger) showProgress(step *Step) {
	elapsed := time.Since(l.StartTime).Round(time.Second)
	
	switch step.Status {
	case StepRunning:
		fmt.Printf("⏳ %s... (%s)\n", step.Name, elapsed)
	case StepComplete:
		duration := step.Duration.Round(time.Millisecond)
		fmt.Printf("✅ %s (%s)\n", step.Name, duration)
	case StepWarning:
		duration := step.Duration.Round(time.Millisecond)
		fmt.Printf("⚠️  %s (%s) - %s\n", step.Name, duration, step.Error.Error())
	case StepError:
		duration := step.Duration.Round(time.Millisecond)
		fmt.Printf("❌ %s (%s) - %s\n", step.Name, duration, step.Error.Error())
	}
}

// PrintSummary prints a final summary with log file location
func (l *Logger) PrintSummary() {
	totalDuration := time.Since(l.StartTime)
	fmt.Printf("\n🎉 Setup completed in %s\n", totalDuration.Round(time.Second))
	fmt.Printf("📄 Full logs: %s\n", l.LogPath)
	fmt.Printf("💡 View logs: steele logs %s\n", l.ProjectName)
}

// GetLatestLogFile returns the path to the latest log file for a project
func GetLatestLogFile(projectName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	logsDir := filepath.Join(homeDir, ".steele", "logs")
	
	// Find all log files for the project
	pattern := filepath.Join(logsDir, fmt.Sprintf("%s_*.log", projectName))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to find log files: %w", err)
	}
	
	if len(matches) == 0 {
		return "", fmt.Errorf("no log files found for project %s", projectName)
	}
	
	// Return the most recent (last in alphabetical order due to timestamp format)
	latest := matches[len(matches)-1]
	return latest, nil
}

// GetAllLogFiles returns all log files for a project
func GetAllLogFiles(projectName string) ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	logsDir := filepath.Join(homeDir, ".steele", "logs")
	
	// Find all log files for the project
	pattern := filepath.Join(logsDir, fmt.Sprintf("%s_*.log", projectName))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find log files: %w", err)
	}
	
	return matches, nil
}