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
	e.txID = ""
	e.inTransaction = false

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

	return lastErr
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
			reconnectErr := c.Connect()
			if reconnectErr != nil {
				return nil, fmt.Errorf("reconnect failed: %w", reconnectErr)
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
