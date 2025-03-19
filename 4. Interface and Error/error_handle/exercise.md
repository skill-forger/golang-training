## Practical Exercises

### Exercise 1: Building a Robust API Client

Create a resilient API client that demonstrates comprehensive error handling practices in Go. This exercise will teach you how to design clear error types, use the error wrapping features, and implement error-based control flow.

Your implementation should include:
1. Custom error types that:
   - Include relevant contextual information (status code, URL, message)
   - Implement the `Error()` method for clear error messages
   - Support unwrapping via the `Unwrap()` method
2. Sentinel (predefined) errors for common error conditions:
   - `ErrNotFound` for missing resources
   - `ErrUnauthorized` for authentication failures
   - `ErrTimeout` for request timeouts
3. An `APIClient` struct with methods for:
   - Making HTTP requests with proper error handling
   - Handling different error cases (timeouts, HTTP error codes)
   - Wrapping underlying errors with context
4. A demonstration showing:
   - Proper error checking with specific error types
   - Using `errors.Is()` to check for sentinel errors
   - Using `errors.As()` to extract information from custom errors
   - Providing appropriate feedback based on error types

```go
// api_client.go
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
)

// Define custom error types
type APIError struct {
    StatusCode int
    URL        string
    Message    string
    Err        error
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error (%d) on %s: %s", e.StatusCode, e.URL, e.Message)
}

func (e *APIError) Unwrap() error {
    return e.Err
}

// Define sentinel errors
var (
    ErrNotFound     = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrTimeout      = errors.New("request timed out")
)

// APIClient for making HTTP requests
type APIClient struct {
    BaseURL    string
    HTTPClient *http.Client
    AuthToken  string
}

// NewAPIClient creates a new client with default settings
func NewAPIClient(baseURL, token string) *APIClient {
    return &APIClient{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 10 * time.Second,
        },
        AuthToken: token,
    }
}

// GetUser fetches a user from the API
func (c *APIClient) GetUser(userID string) (map[string]interface{}, error) {
    url := fmt.Sprintf("%s/users/%s", c.BaseURL, userID)
    
    // Create request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    // Add authorization if available
    if c.AuthToken != "" {
        req.Header.Set("Authorization", "Bearer "+c.AuthToken)
    }
    
    // Make the request
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        // Handle timeout specifically
        if errors.Is(err, http.ErrHandlerTimeout) {
            return nil, ErrTimeout
        }
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    // Read response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    // Handle different status codes
    switch resp.StatusCode {
    case http.StatusOK:
        // Success - parse the JSON
        var user map[string]interface{}
        if err := json.Unmarshal(body, &user); err != nil {
            return nil, fmt.Errorf("failed to parse response: %w", err)
        }
        return user, nil
        
    case http.StatusNotFound:
        return nil, &APIError{
            StatusCode: resp.StatusCode,
            URL:        url,
            Message:    "User not found",
            Err:        ErrNotFound,
        }
        
    case http.StatusUnauthorized:
        return nil, &APIError{
            StatusCode: resp.StatusCode,
            URL:        url,
            Message:    "Invalid or expired token",
            Err:        ErrUnauthorized,
        }
        
    default:
        // Generic error for other status codes
        return nil, &APIError{
            StatusCode: resp.StatusCode,
            URL:        url,
            Message:    fmt.Sprintf("API returned status %d", resp.StatusCode),
            Err:        errors.New("unexpected API response"),
        }
    }
}

func main() {
    // Create a client
    client := NewAPIClient("https://api.example.com", "valid-token")
    
    // Make a request
    user, err := client.GetUser("123")
    
    // Handle errors with appropriate type checks
    if err != nil {
        var apiErr *APIError
        
        switch {
        case errors.Is(err, ErrNotFound):
            fmt.Println("User not found")
            
        case errors.Is(err, ErrUnauthorized):
            fmt.Println("Please log in again")
            
        case errors.Is(err, ErrTimeout):
            fmt.Println("Request timed out, please try again")
            
        case errors.As(err, &apiErr):
            fmt.Printf("API error (%d): %s\n", apiErr.StatusCode, apiErr.Message)
            
        default:
            fmt.Printf("Unexpected error: %v\n", err)
        }
        
        return
    }
    
    // Process user data
    fmt.Printf("User: %v\n", user)
}
```

### Exercise 2: Database Connection with Error Recovery

Develop a database connection manager that handles various error scenarios and implements automatic recovery strategies. This exercise shows how to use errors to make robust systems that can recover from failures.

Your implementation should include:
1. Custom error types for different database issues:
   - Connection errors
   - Query execution errors
   - Transaction errors
2. A `DBConnector` that manages database connections with:
   - Connection retry logic with exponential backoff
   - Transaction handling with proper rollback on errors
   - Query execution with timeout handling
3. A `QueryExecutor` interface with implementations for:
   - Basic queries
   - Prepared statements
   - Transactions
4. A demonstration showing:
   - Handling temporary network failures
   - Automatic reconnection after errors
   - Transaction rollback on errors
   - Proper resource cleanup

```go
// db_connector.go
package main

import (
    "errors"
    "fmt"
    "math"
    "math/rand"
    "time"
)

// Define custom error types
type DBError struct {
    Operation string
    Message   string
    Err       error
}

func (e *DBError) Error() string {
    return fmt.Sprintf("database error during %s: %s", e.Operation, e.Message)
}

func (e *DBError) Unwrap() error {
    return e.Err
}

// Specific database error types
type ConnectionError struct {
    DBError
    ConnectionString string
}

func (e *ConnectionError) Error() string {
    return fmt.Sprintf("failed to connect to database at %s: %s", 
        e.ConnectionString, e.Message)
}

type QueryError struct {
    DBError
    Query string
}

func (e *QueryError) Error() string {
    return fmt.Sprintf("query failed [%s]: %s", e.Query, e.Message)
}

type TransactionError struct {
    DBError
    TxID string
}

func (e *TransactionError) Error() string {
    return fmt.Sprintf("transaction %s failed: %s", e.TxID, e.Message)
}

// Sentinel errors
var (
    ErrConnectionFailed  = errors.New("database connection failed")
    ErrQueryFailed       = errors.New("query execution failed")
    ErrTransactionFailed = errors.New("transaction failed")
    ErrTimeout           = errors.New("operation timed out")
)

// QueryExecutor defines methods for database operations
type QueryExecutor interface {
    Execute(query string, args ...interface{}) (interface{}, error)
}

// BasicExecutor implements simple query execution
type BasicExecutor struct {
    connected bool
    failRate  float64 // Simulate random failures
}

func (e *BasicExecutor) Execute(query string, args ...interface{}) (interface{}, error) {
    if !e.connected {
        return nil, &QueryError{
            DBError: DBError{
                Operation: "execute",
                Message:   "not connected to database",
                Err:       ErrConnectionFailed,
            },
            Query: query,
        }
    }
    
    // Simulate random failures
    if rand.Float64() < e.failRate {
        return nil, &QueryError{
            DBError: DBError{
                Operation: "execute",
                Message:   "random failure occurred",
                Err:       ErrQueryFailed,
            },
            Query: query,
        }
    }
    
    // Simulate successful execution
    return fmt.Sprintf("Executed: %s with args %v", query, args), nil
}

// TransactionExecutor handles database transactions
type TransactionExecutor struct {
    BasicExecutor
    inTransaction bool
    txID          string
}

func (e *TransactionExecutor) BeginTransaction() error {
    if !e.connected {
        return &ConnectionError{
            DBError: DBError{
                Operation: "begin transaction",
                Message:   "not connected to database",
                Err:       ErrConnectionFailed,
            },
            ConnectionString: "db:3306",
        }
    }
    
    if e.inTransaction {
        return &TransactionError{
            DBError: DBError{
                Operation: "begin transaction",
                Message:   "already in a transaction",
                Err:       ErrTransactionFailed,
            },
            TxID: e.txID,
        }
    }
    
    // Start a new transaction
    e.inTransaction = true
    e.txID = fmt.Sprintf("tx-%d", time.Now().UnixNano())
    return nil
}

func (e *TransactionExecutor) Commit() error {
    if !e.inTransaction {
        return &TransactionError{
            DBError: DBError{
                Operation: "commit",
                Message:   "no active transaction",
                Err:       ErrTransactionFailed,
            },
            TxID: "",
        }
    }
    
    // Simulate random commit failures
    if rand.Float64() < e.failRate {
        return &TransactionError{
            DBError: DBError{
                Operation: "commit",
                Message:   "failed to commit transaction",
                Err:       ErrTransactionFailed,
            },
            TxID: e.txID,
        }
    }
    
    // Commit successful
    e.inTransaction = false
    txID := e.txID
    e.txID = ""
    
    return nil
}

func (e *TransactionExecutor) Rollback() error {
    if !e.inTransaction {
        return &TransactionError{
            DBError: DBError{
                Operation: "rollback",
                Message:   "no active transaction",
                Err:       ErrTransactionFailed,
            },
            TxID: "",
        }
    }
    
    // Rollback is almost always successful
    e.inTransaction = false
    e.txID = ""
    
    return nil
}

// DBConnector manages database connections
type DBConnector struct {
    dsn          string
    connected    bool
    executor     QueryExecutor
    maxRetries   int
    retryBackoff time.Duration
}

func NewDBConnector(dsn string) *DBConnector {
    return &DBConnector{
        dsn:          dsn,
        connected:    false,
        maxRetries:   5,
        retryBackoff: 100 * time.Millisecond,
    }
}

func (c *DBConnector) Connect() error {
    // Try to connect with retries
    var lastErr error
    
    for i := 0; i < c.maxRetries; i++ {
        // Simulate connection attempt
        success := rand.Float64() > 0.3 // 70% success rate
        
        if success {
            c.connected = true
            
            // Create executor with 20% failure rate
            basicExec := &BasicExecutor{
                connected: true,
                failRate:  0.2,
            }
            
            c.executor = &TransactionExecutor{
                BasicExecutor: *basicExec,
                inTransaction: false,
            }
            
            return nil
        }
        
        // Connection failed, calculate backoff
        backoff := c.retryBackoff * time.Duration(math.Pow(2, float64(i)))
        
        lastErr = &ConnectionError{
            DBError: DBError{
                Operation: "connect",
                Message:   fmt.Sprintf("attempt %d failed, retrying in %v", i+1, backoff),
                Err:       ErrConnectionFailed,
            },
            ConnectionString: c.dsn,
        }
        
        fmt.Printf("Connection attempt %d failed, retrying in %v\n", i+1, backoff)
        time.Sleep(backoff)
    }
    
    return &ConnectionError{
        DBError: DBError{
            Operation: "connect",
            Message:   fmt.Sprintf("failed after %d attempts", c.maxRetries),
            Err:       ErrConnectionFailed,
        },
        ConnectionString: c.dsn,
    }
}

func (c *DBConnector) Disconnect() error {
    if !c.connected {
        return nil // Already disconnected
    }
    
    c.connected = false
    c.executor = nil
    return nil
}

func (c *DBConnector) Execute(query string, args ...interface{}) (interface{}, error) {
    if !c.connected {
        err := c.Connect()
        if err != nil {
            return nil, fmt.Errorf("auto-connect failed: %w", err)
        }
    }
    
    result, err := c.executor.Execute(query, args...)
    
    // Handle specific errors
    if err != nil {
        var connErr *ConnectionError
        if errors.As(err, &connErr) {
            // Try to reconnect
            c.connected = false
            reconnErr := c.Connect()
            if reconnErr != nil {
                return nil, fmt.Errorf("reconnect failed: %w", reconnErr)
            }
            
            // Retry the query
            return c.executor.Execute(query, args...)
        }
        
        return nil, err
    }
    
    return result, nil
}

func (c *DBConnector) ExecuteTransaction(queries []string) ([]interface{}, error) {
    if !c.connected {
        err := c.Connect()
        if err != nil {
            return nil, fmt.Errorf("auto-connect failed: %w", err)
        }
    }
    
    // Get transaction executor
    txExecutor, ok := c.executor.(*TransactionExecutor)
    if !ok {
        return nil, fmt.Errorf("executor does not support transactions")
    }
    
    // Begin transaction
    err := txExecutor.BeginTransaction()
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    
    // Execute all queries
    results := make([]interface{}, 0, len(queries))
    
    for _, query := range queries {
        result, err := txExecutor.Execute(query)
        if err != nil {
            // Transaction failed, roll back
            rollbackErr := txExecutor.Rollback()
            if rollbackErr != nil {
                return nil, fmt.Errorf("query failed (%w) and rollback failed (%v)", 
                    err, rollbackErr)
            }
            
            return nil, fmt.Errorf("query failed, transaction rolled back: %w", err)
        }
        
        results = append(results, result)
    }
    
    // Commit transaction
    err = txExecutor.Commit()
    if err != nil {
        // Commit failed, roll back
        rollbackErr := txExecutor.Rollback()
        if rollbackErr != nil {
            return nil, fmt.Errorf("commit failed (%w) and rollback failed (%v)", 
                err, rollbackErr)
        }
        
        return nil, fmt.Errorf("commit failed, transaction rolled back: %w", err)
    }
    
    return results, nil
}

func main() {
    // Seed the random number generator
    rand.Seed(time.Now().UnixNano())
    
    // Create a database connector
    db := NewDBConnector("mysql://localhost:3306/testdb")
    
    // Connect to the database
    err := db.Connect()
    if err != nil {
        var connErr *ConnectionError
        if errors.As(err, &connErr) {
            fmt.Printf("Connection error: %v\n", connErr)
        } else {
            fmt.Printf("Error: %v\n", err)
        }
        return
    }
    
    fmt.Println("Connected to database successfully")
    
    // Execute a simple query
    fmt.Println("\n--- Executing single query ---")
    result, err := db.Execute("SELECT * FROM users WHERE id = ?", 1)
    if err != nil {
        var queryErr *QueryError
        if errors.As(err, &queryErr) {
            fmt.Printf("Query error: %v\n", queryErr)
        } else {
            fmt.Printf("Error: %v\n", err)
        }
    } else {
        fmt.Printf("Result: %v\n", result)
    }
    
    // Execute a transaction
    fmt.Println("\n--- Executing transaction ---")
    queries := []string{
        "INSERT INTO users (name, email) VALUES ('John', 'john@example.com')",
        "UPDATE users SET status = 'active' WHERE email = 'john@example.com'",
    }
    
    results, err := db.ExecuteTransaction(queries)
    if err != nil {
        var txErr *TransactionError
        if errors.As(err, &txErr) {
            fmt.Printf("Transaction error: %v\n", txErr)
        } else if errors.Is(err, ErrTransactionFailed) {
            fmt.Printf("Generic transaction failure: %v\n", err)
        } else {
            fmt.Printf("Error: %v\n", err)
        }
    } else {
        fmt.Println("Transaction completed successfully:")
        for i, result := range results {
            fmt.Printf("Query %d result: %v\n", i+1, result)
        }
    }
    
    // Disconnect
    db.Disconnect()
    fmt.Println("\nDisconnected from database")
}
```

### Exercise 3: File Processing with Error Logging

Build a file processing system that demonstrates advanced error handling techniques, including logging, error wrapping, and recovery mechanisms. This exercise shows how to handle I/O errors and create user-friendly error messages.

Your implementation should include:
1. A hierarchical error handling system with:
   - Base error types for I/O operations
   - Specialized errors for different file processing stages
   - Error aggregation for batch operations
2. A structured logging system that:
   - Records errors with appropriate context
   - Categorizes errors by severity
   - Provides detailed debugging information
3. Recovery mechanisms that:
   - Skip problematic files and continue processing others
   - Attempt alternative processing methods when primary methods fail
   - Clean up resources even when errors occur
4. A demonstration that processes multiple files and shows how the system handles various error conditions

```go
// file_processor.go
package main

import (
    "bufio"
    "errors"
    "fmt"
    "io"
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

// Convenience logger methods
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

// Define custom error types for file processing
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
```
