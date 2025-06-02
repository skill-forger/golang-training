# Module 09: HTTP Servers in Go - Building Web Applications

## Introduction to HTTP Servers in Go

Go provides powerful, yet straightforward tools for building web applications through its standard library. The `net/http` package offers everything needed to create robust HTTP servers without requiring external dependencies, embodying Go's philosophy of simplicity and efficiency.

### The HTTP Server Landscape

Before exploring Go's HTTP server implementation, let's understand the fundamentals:

1. **HTTP Protocol Basics**
    - **Request-Response Model**: How clients communicate with servers
    - **HTTP Methods**: GET, POST, PUT, DELETE, etc., and their intended uses
    - **Status Codes**: Standardized responses (200 OK, 404 Not Found, etc.)

2. **Why HTTP Servers in Go Matter**
    - Built-in concurrency for handling multiple requests
    - Minimal memory footprint and high performance
    - Simple API that encourages clean code organization

### Creating a Basic HTTP Server

The standard library's `net/http` package makes it incredibly easy to create web servers with just a few lines of code:

```go
// Basic HTTP server example
func main() {
    // Define a handler function
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, you've requested: %s", r.URL.Path)
    })

    // Start the server on port 8080
    fmt.Println("Starting server at port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
```

#### HTTP Server Components
- **Handler Functions**: Process incoming requests
- **ServeMux (Router)**: Routes requests to appropriate handlers
- **Server**: Listens for connections and dispatches to handlers
- **ResponseWriter**: Constructs HTTP responses
- **Request**: Contains data about the incoming HTTP request

### Handling HTTP Requests

Go's HTTP server revolves around the `http.Handler` interface, which defines how to process incoming requests:

```go
// Custom handler implementation
type customHandler struct{}

func (h customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Custom handler serving: %s", r.URL.Path)
}

func main() {
    handler := customHandler{}
    
    // Register our custom handler for all routes
    http.Handle("/", handler)
    
    // Start the server
    http.ListenAndServe(":8080", nil)
}
```

#### Working with Request Data
```go
func formHandler(w http.ResponseWriter, r *http.Request) {
    // Parse form data
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }
    
    // Access form values
    name := r.FormValue("name")
    
    // Write response
    fmt.Fprintf(w, "Form submission successful. Name = %s", name)
}

func main() {
    http.HandleFunc("/form", formHandler)
    http.ListenAndServe(":8080", nil)
}
```

### Custom Router with ServeMux

The `http.ServeMux` provides a flexible routing mechanism for your web applications:

```go
func main() {
    // Create a new router
    mux := http.NewServeMux()
    
    // Register route handlers
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/about", aboutHandler)
    mux.HandleFunc("/api/", apiHandler)
    
    // Start server with custom router
    http.ListenAndServe(":8080", mux)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    fmt.Fprintf(w, "Welcome to the home page!")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "About Us")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    path := strings.TrimPrefix(r.URL.Path, "/api/")
    fmt.Fprintf(w, "API request for: %s", path)
}
```

### Serving Static Files

Go makes it easy to serve static assets like images, CSS, and JavaScript:

```go
func main() {
    // Create file server handler
    fs := http.FileServer(http.Dir("./static"))
    
    // Register the handler
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    
    // Register other routes
    http.HandleFunc("/", indexHandler)
    
    // Start the server
    http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    // Serve an HTML page that references static files
    html := `
    <!DOCTYPE html>
    <html>
        <head>
            <title>Go Web Server</title>
            <link rel="stylesheet" href="/static/style.css">
        </head>
        <body>
            <h1>Welcome to Go Web Development</h1>
            <img src="/static/gopher.png" alt="Go Gopher">
            <script src="/static/script.js"></script>
        </body>
    </html>
    `
    fmt.Fprint(w, html)
}
```

### HTTP Middleware

Middleware functions in Go allow you to wrap handlers with common functionality:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Pre-processing logic
        startTime := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)
        
        // Call the next handler
        next.ServeHTTP(w, r)
        
        // Post-processing logic
        log.Printf("Completed %s in %v", r.URL.Path, time.Since(startTime))
    })
}

func main() {
    // Create a handler
    handler := http.HandlerFunc(homeHandler)
    
    // Wrap it with middleware
    wrappedHandler := loggingMiddleware(handler)
    
    // Register the wrapped handler
    http.Handle("/", wrappedHandler)
    http.ListenAndServe(":8080", nil)
}
```

### JSON Handling

Working with JSON is common in web applications and APIs:

```go
type User struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
}

func apiUserHandler(w http.ResponseWriter, r *http.Request) {
    // Create sample user
    user := User{
        ID:       1,
        Username: "gopher",
        Email:    "gopher@example.com",
    }
    
    // Set content type header
    w.Header().Set("Content-Type", "application/json")
    
    // Encode and send JSON response
    if err := json.NewEncoder(w).Encode(user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func apiCreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Only accept POST requests
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Parse the request body
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Process the user (in a real app, save to database, etc.)
    
    // Return created status code
    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "User %s created successfully", user.Username)
}
```

### HTTPS Server

Securing your web application with TLS/SSL:

```go
func main() {
    // Configure TLS server
    server := &http.Server{
        Addr:         ":8443",
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    
    // Register handlers
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Secure HTTPS server")
    })
    
    // Start HTTPS server
    log.Println("Starting HTTPS server on :8443")
    log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
```

### Server Configuration and Timeouts

Properly configuring your HTTP server for production:

```go
func main() {
    // Create a custom server with configurations
    server := &http.Server{
        Addr:         ":8080",
        Handler:      http.DefaultServeMux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }
    
    // Register handlers
    http.HandleFunc("/", homeHandler)
    
    // Graceful shutdown
    go func() {
        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, os.Interrupt)
        <-sigint
        
        // Create shutdown context
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        
        // Shutdown the server
        if err := server.Shutdown(ctx); err != nil {
            log.Printf("HTTP server shutdown error: %v", err)
        }
    }()
    
    // Start server
    log.Printf("Starting HTTP server on %s", server.Addr)
    if err := server.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("HTTP server error: %v", err)
    }
    log.Println("Server gracefully stopped")
}
```

### Common HTTP Server Challenges

1. **Performance Optimization**
    - Connection pooling
    - Response compression
    - Efficient routing

2. **Security Concerns**
    - Input validation
    - CSRF protection
    - Rate limiting

3. **Error Handling**
    - Consistent error responses
    - Logging and monitoring
    - Graceful failure modes

### Best Practices

1. Use the appropriate HTTP status codes
2. Implement proper request timeouts
3. Structure your API paths consistently
4. Log requests and errors for debugging
5. Handle all errors explicitly

### Learning Challenges

1. Create a RESTful API with CRUD operations
2. Build a file upload server
3. Implement user authentication
4. Create a reverse proxy server
5. Build a WebSocket server for real-time communication

### Recommended Resources
- "Web Development with Go" by Shiju Varghese
- Go's official documentation for net/http package
- "Building Web Applications with Go" on the Go website

### Reflection Questions

1. How does Go's concurrency model benefit HTTP server performance?
2. What are the trade-offs between using Go's standard library versus frameworks like Gin or Echo?
3. How would you design a scalable web service architecture using Go's HTTP packages?

**HTTP Server Mastery: Build Powerful Web Applications with Go's Standard Library!** 🌐