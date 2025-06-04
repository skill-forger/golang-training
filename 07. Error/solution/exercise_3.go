package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogLevel defines different logging severity levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}[l]
}

// Logger provides structured logging functionality
type Logger struct {
	Level   LogLevel
	LogFile *os.File
}

// NewLogger creates a new logger with the specified minimum level
func NewLogger(level LogLevel, logPath string) (*Logger, error) {
	var logFile *os.File
	var err error

	if logPath != "" {
		logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
	}

	return &Logger{
		Level:   level,
		LogFile: logFile,
	}, nil
}

// Log writes a log entry with the given level and message
func (l *Logger) Log(level LogLevel, format string, args ...interface{}) {
	if level < l.Level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logEntry := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)

	// Write to console
	fmt.Print(logEntry)

	// Write to file if available
	if l.LogFile != nil {
		l.LogFile.WriteString(logEntry)
	}
}

// Close closes the log file if it's open
func (l *Logger) Close() error {
	if l.LogFile != nil {
		return l.LogFile.Close()
	}
	return nil
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.Log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.Log(INFO, format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.Log(WARNING, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.Log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.Log(FATAL, format, args...)
}

// FileError defines custom error types for file processing
type FileError struct {
	Path    string
	Op      string
	Message string
	Err     error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("file error [%s] on %s: %s", e.Op, e.Path, e.Message)
}

func (e *FileError) Unwrap() error {
	return e.Err
}

// Specific file error types
type ReadError struct {
	FileError
	Line int
}

func (e *ReadError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("read error at line %d in %s: %s", e.Line, e.Path, e.Message)
	}
	return fmt.Sprintf("read error in %s: %s", e.Path, e.Message)
}

type ParseError struct {
	FileError
	Line    int
	Content string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d in %s: %s (content: %q)",
		e.Line, e.Path, e.Message, e.Content)
}

type WriteError struct {
	FileError
}

func (e *WriteError) Error() string {
	return fmt.Sprintf("write error in %s: %s", e.Path, e.Message)
}

// BatchError collects multiple errors from batch operations
type BatchError struct {
	Errors []error
}

func (e *BatchError) Error() string {
	if len(e.Errors) == 1 {
		return fmt.Sprintf("batch operation failed with 1 error: %v", e.Errors[0])
	}
	return fmt.Sprintf("batch operation failed with %d errors", len(e.Errors))
}

func (e *BatchError) AddError(err error) {
	e.Errors = append(e.Errors, err)
}

func (e *BatchError) HasErrors() bool {
	return len(e.Errors) > 0
}

// Define sentinel errors
var (
	ErrFileNotFound = errors.New("file not found")
	ErrPermission   = errors.New("permission denied")
	ErrFormat       = errors.New("invalid file format")
	ErrEmpty        = errors.New("file is empty")
)

// FileProcessor handles file processing operations
type FileProcessor struct {
	Logger *Logger
}

// NewFileProcessor creates a new file processor
func NewFileProcessor(logger *Logger) *FileProcessor {
	return &FileProcessor{
		Logger: logger,
	}
}

// ProcessFile reads a file and performs line-by-line processing
func (p *FileProcessor) ProcessFile(path string) error {
	p.Logger.Info("Processing file: %s", path)

	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &FileError{
			Path:    path,
			Op:      "check",
			Message: "file does not exist",
			Err:     ErrFileNotFound,
		}
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		if os.IsPermission(err) {
			return &FileError{
				Path:    path,
				Op:      "open",
				Message: "permission denied",
				Err:     ErrPermission,
			}
		}
		return &FileError{
			Path:    path,
			Op:      "open",
			Message: err.Error(),
			Err:     err,
		}
	}
	defer file.Close()

	// Prepare output file
	outputPath := filepath.Join(
		filepath.Dir(path),
		fmt.Sprintf("processed_%s", filepath.Base(path)),
	)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return &WriteError{
			FileError: FileError{
				Path:    outputPath,
				Op:      "create",
				Message: err.Error(),
				Err:     err,
			},
		}
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	// Process the file line by line
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if line == "" {
			p.Logger.Debug("Skipping empty line %d", lineNum)
			continue
		}

		// Process the line
		processed, err := p.processLine(path, lineNum, line)
		if err != nil {
			// Log the error but continue processing
			p.Logger.Warning("Error processing line %d: %v", lineNum, err)

			// Write error comment to output
			fmt.Fprintf(writer, "# ERROR Line %d: %s\n", lineNum, err.Error())
			continue
		}

		// Write the processed line
		if _, err := fmt.Fprintln(writer, processed); err != nil {
			return &WriteError{
				FileError: FileError{
					Path:    outputPath,
					Op:      "write",
					Message: err.Error(),
					Err:     err,
				},
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return &ReadError{
			FileError: FileError{
				Path:    path,
				Op:      "scan",
				Message: err.Error(),
				Err:     err,
			},
			Line: lineNum,
		}
	}

	if lineNum == 0 {
		return &FileError{
			Path:    path,
			Op:      "process",
			Message: "file is empty",
			Err:     ErrEmpty,
		}
	}

	p.Logger.Info("Successfully processed %s, wrote output to %s", path, outputPath)
	return nil
}

// processLine handles a single line of text
func (p *FileProcessor) processLine(path string, lineNum int, line string) (string, error) {
	// Example processing: Convert CSV to pipe-delimited format
	if !strings.Contains(line, ",") {
		return "", &ParseError{
			FileError: FileError{
				Path:    path,
				Op:      "parse",
				Message: "line does not contain delimiters",
				Err:     ErrFormat,
			},
			Line:    lineNum,
			Content: line,
		}
	}

	// Split by comma and rejoin with pipe
	fields := strings.Split(line, ",")

	// Trim whitespace from each field
	for i, field := range fields {
		fields[i] = strings.TrimSpace(field)
	}

	return strings.Join(fields, "|"), nil
}

// ProcessFiles processes multiple files and collects errors
func (p *FileProcessor) ProcessFiles(paths []string) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	batchErr := &BatchError{}

	for _, path := range paths {
		err := p.ProcessFile(path)
		if err != nil {
			p.Logger.Error("Failed to process %s: %v", path, err)
			batchErr.AddError(err)
		}
	}

	if batchErr.HasErrors() {
		return batchErr
	}

	return nil
}

func main() {
	// Create a logger
	logger, err := NewLogger(DEBUG, "file_processor.log")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	processor := NewFileProcessor(logger)

	// Create some test files
	testDir := "test_files"
	os.Mkdir(testDir, 0755)

	files := map[string][]string{
		"valid.csv": {
			"Name, Age, City",
			"John, 30, New York",
			"Alice, 25, London",
			"Bob, 40, Paris",
		},
		"empty.csv": {},
		"invalid.txt": {
			"This is not a CSV file",
			"It has no commas at all",
		},
		"partially_valid.csv": {
			"Product, Price, Quantity",
			"Apple, 1.20, 10",
			"Invalid line with no commas",
			"Orange, 0.80, 5",
		},
	}

	var filePaths []string

	for name, content := range files {
		path := filepath.Join(testDir, name)
		filePaths = append(filePaths, path)

		// Create the file
		file, err := os.Create(path)
		if err != nil {
			logger.Error("Failed to create test file %s: %v", path, err)
			continue
		}

		// Write content
		for _, line := range content {
			fmt.Fprintln(file, line)
		}

		file.Close()
	}

	// Process all files
	err = processor.ProcessFiles(filePaths)

	// Handle batch errors
	if err != nil {
		var batchErr *BatchError
		if errors.As(err, &batchErr) {
			logger.Error("Batch processing completed with %d errors", len(batchErr.Errors))

			for i, err := range batchErr.Errors {
				var fileErr *FileError
				var readErr *ReadError
				var parseErr *ParseError
				var writeErr *WriteError

				switch {
				case errors.As(err, &parseErr):
					logger.Error("Error %d: Parse error at line %d in %s: %s",
						i+1, parseErr.Line, parseErr.Path, parseErr.Message)

				case errors.As(err, &readErr):
					logger.Error("Error %d: Read error at line %d in %s: %s",
						i+1, readErr.Line, readErr.Path, readErr.Message)

				case errors.As(err, &writeErr):
					logger.Error("Error %d: Write error in %s: %s",
						i+1, writeErr.Path, writeErr.Message)

				case errors.As(err, &fileErr):
					logger.Error("Error %d: File error [%s] on %s: %s",
						i+1, fileErr.Op, fileErr.Path, fileErr.Message)

				default:
					logger.Error("Error %d: %v", i+1, err)
				}
			}
		} else {
			logger.Error("Error processing files: %v", err)
		}
	} else {
		logger.Info("All files processed successfully")
	}
}
