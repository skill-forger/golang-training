## Practical Exercises

### Exercise 1: Simple RESTful API

Build a simple RESTful API to manage a collection of books:

```go
// book_api.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Book represents a book entity
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

// BookStore manages the collection of books
type BookStore struct {
	books  []Book
	nextID int
}

// NewBookStore creates a new book store with some initial data
func NewBookStore() *BookStore {
	return &BookStore{
		books: []Book{
			{ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan & Brian Kernighan", Year: 2015},
			{ID: 2, Title: "Go in Action", Author: "William Kennedy", Year: 2016},
		},
		nextID: 3,
	}
}

func main() {
	store := NewBookStore()
	
	// Define handlers
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// GET /books - Return all books
			handleGetBooks(w, store)
		case http.MethodPost:
			// POST /books - Create a new book
			handleCreateBook(w, r, store)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	
	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		// Extract book ID from URL
		idStr := strings.TrimPrefix(r.URL.Path, "/books/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid book ID", http.StatusBadRequest)
			return
		}
		
		switch r.Method {
		case http.MethodGet:
			// GET /books/{id} - Get a specific book
			handleGetBook(w, id, store)
		case http.MethodPut:
			// PUT /books/{id} - Update a specific book
			handleUpdateBook(w, r, id, store)
		case http.MethodDelete:
			// DELETE /books/{id} - Delete a specific book
			handleDeleteBook(w, id, store)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	
	// Start server
	fmt.Println("Starting book server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler functions

func handleGetBooks(w http.ResponseWriter, store *BookStore) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.books)
}

func handleCreateBook(w http.ResponseWriter, r *http.Request, store *BookStore) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Assign a new ID
	book.ID = store.nextID
	store.nextID++
	
	// Add to collection
	store.books = append(store.books, book)
	
	// Return the created book
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func handleGetBook(w http.ResponseWriter, id int, store *BookStore) {
	for _, book := range store.books {
		if book.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	
	http.Error(w, "Book not found", http.StatusNotFound)
}

func handleUpdateBook(w http.ResponseWriter, r *http.Request, id int, store *BookStore) {
	var updatedBook Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	for i, book := range store.books {
		if book.ID == id {
			// Preserve the book ID
			updatedBook.ID = id
			store.books[i] = updatedBook
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}
	
	http.Error(w, "Book not found", http.StatusNotFound)
}

func handleDeleteBook(w http.ResponseWriter, id int, store *BookStore) {
	for i, book := range store.books {
		if book.ID == id {
			// Remove the book
			store.books = append(store.books[:i], store.books[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	
	http.Error(w, "Book not found", http.StatusNotFound)
}
```

### Exercise 2: File Server with Custom Handler

Create a file server with a custom middleware for logging:

```go
// file_server.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// LoggingMiddleware adds logging to each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		startTime := time.Now()
		fmt.Printf("[%s] %s %s\n", startTime.Format(time.RFC822), r.Method, r.URL.Path)
		
		// Call the next handler
		next.ServeHTTP(w, r)
		
		// Log the time taken
		duration := time.Since(startTime)
		fmt.Printf("Request completed in %v\n", duration)
	})
}

// NotFoundHandler handles 404 errors
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`
		<html>
			<head><title>Not Found</title></head>
			<body>
				<h1>404 - Page Not Found</h1>
				<p>The page you requested does not exist.</p>
				<a href="/">Go to Home</a>
			</body>
		</html>
	`))
}

func main() {
	// Create the directory if it doesn't exist
	os.MkdirAll("./static", 0755)
	
	// Create a sample index.html file
	indexContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Go File Server</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
				h1 { color: #333; }
			</style>
		</head>
		<body>
			<h1>Welcome to the Go File Server</h1>
			<p>This is a simple file server built with Go's standard library.</p>
		</body>
		</html>
	`
	os.WriteFile("./static/index.html", []byte(indexContent), 0644)
	
	// Create a file server handler
	fileServer := http.FileServer(http.Dir("./static"))
	
	// Register handlers
	mux := http.NewServeMux()
	
	// Add the file server with the logging middleware
	mux.Handle("/", LoggingMiddleware(fileServer))
	
	// Custom not found handler
	mux.HandleFunc("/notfound", NotFoundHandler)
	
	// Start the server
	fmt.Println("Starting file server on :8080...")
	fmt.Println("Files are served from the ./static directory")
	fmt.Println("Use Ctrl+C to stop the server")
	
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### Exercise 3: HTTP Client and Concurrent Requests

Build an HTTP client that makes concurrent requests to different APIs:

```go
// concurrent_client.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// ApiResponse represents a generic API response
type ApiResponse struct {
	Source  string
	Data    map[string]interface{}
	Error   error
	Latency time.Duration
}

// FetchAPI makes an HTTP request to the given URL and returns the response
func FetchAPI(url string, source string) ApiResponse {
	startTime := time.Now()
	resp, err := http.Get(url)
	latency := time.Since(startTime)
	
	if err != nil {
		return ApiResponse{
			Source:  source,
			Error:   err,
			Latency: latency,
		}
	}
	defer resp.Body.Close()
	
	var data map[string]interface{}
	
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return ApiResponse{
				Source:  source,
				Error:   fmt.Errorf("failed to read response body: %w", err),
				Latency: latency,
			}
		}
		
		if err := json.Unmarshal(body, &data); err != nil {
			return ApiResponse{
				Source:  source,
				Error:   fmt.Errorf("failed to parse JSON: %w", err),
				Latency: latency,
			}
		}
	} else {
		return ApiResponse{
			Source:  source,
			Error:   fmt.Errorf("API returned status code %d", resp.StatusCode),
			Latency: latency,
		}
	}
	
	return ApiResponse{
		Source:  source,
		Data:    data,
		Latency: latency,
	}
}

func main() {
	// List of APIs to fetch (using httpbin for demonstration)
	apis := []struct {
		URL    string
		Source string
	}{
		{"https://httpbin.org/get", "HTTPBin Get"},
		{"https://httpbin.org/ip", "IP Info"},
		{"https://httpbin.org/user-agent", "User Agent"},
		{"https://httpbin.org/headers", "Headers"},
		{"https://httpbin.org/delay/2", "Delayed Response"}, // This one will take longer
	}
	
	// Channel to collect responses
	responses := make(chan ApiResponse, len(apis))
	
	// WaitGroup to wait for all goroutines
	var wg sync.WaitGroup
	wg.Add(len(apis))
	
	// Set a custom timeout for the HTTP client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	// Override the default HTTP client
	http.DefaultClient = client
	
	fmt.Println("Making concurrent API requests...")
	
	// Make requests concurrently
	for _, api := range apis {
		go func(url, source string) {
			defer wg.Done()
			responses <- FetchAPI(url, source)
		}(api.URL, api.Source)
	}
	
	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(responses)
	}()
	
	// Process the results
	for resp := range responses {
		if resp.Error != nil {
			fmt.Printf("[%s] Error: %v (took %v)\n", resp.Source, resp.Error, resp.Latency)
		} else {
			fmt.Printf("[%s] Success (took %v)\n", resp.Source, resp.Latency)
			// Print a sample of the data
			for k, v := range resp.Data {
				fmt.Printf("  - %s: %v\n", k, v)
				break // Just show one item for brevity
			}
		}
	}
	
	fmt.Println("All requests completed!")
}
```
