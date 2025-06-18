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
