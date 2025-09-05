# Module 12: Gin Gonic

## Table of Contents

<ol>
    <li><a href="#objectives">Objectives</a></li>
    <li><a href="#overview">Overview</a></li>
    <li><a href="#what-is-a-server">What is a Server?</a></li>
    <li><a href="#http-servers">HTTP Servers</a></li>
    <li><a href="#http-server-in-go">HTTP Server in Go</a></li>
    <li><a href="#common-mistakes">Common Mistakes</a></li>
    <li><a href="#best-practices">Best Practices</a></li>
    <li><a href="#practice-exercises">Practice Exercises</a></li>
</ol>

## Objectives

By the end of this module, you will:

- Understand the Core Concepts of a Gin Gonic
- 
- Grasp the HTTP Request-Response Cycle
- Utilize the `net/http` in Go to host a simple HTTP server
- Implement basic use cases and handlers using `net/http` package
- Gain Practical Awareness for Backend Development

## Overview

In this module, the focus will be on the foundational principles of HTTP Servers, which are integral to nearly all
internet interactions. This module aims to equip freshers and interns with a clear conceptual understanding of how
these servers operate, gain insight into the core mechanisms of web communication, including the HTTP
protocol, the structure of requests and responses, and the server's role in processing these exchanges, thereby building
a solid foundation for backend development journey.


## Introduction to Gin Framework

While Go's standard library provides robust HTTP server capabilities, the Gin framework has emerged as one of the most popular web frameworks in the Go ecosystem. Gin offers a more feature-rich, expressive, and performance-focused approach to building web applications while maintaining Go's simplicity and efficiency.

### The Gin Framework Landscape

Before diving into Gin's specifics, let's understand why it has become a leading choice:

1. **Gin vs Standard Library**
    - **Performance**: Gin is built for speed with radix tree-based routing
    - **Developer Experience**: More intuitive API design and middleware system
    - **Feature Set**: Built-in support for JSON validation, error management, and more

2. **Why Gin Matters for Modern Applications**
    - Simplified API development process
    - Robust middleware ecosystem
    - Excellent performance characteristics
    - Clean code organization through groups and routes

### Getting Started with Gin

To begin working with Gin, you'll need to install it first:

```go
// First, install Gin using Go modules
// In your terminal:
// go get -u github.com/gin-gonic/gin

// Basic Gin server example
package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    // Create a default Gin router
    r := gin.Default()

    // Define a route
    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Welcome to Gin Framework!",
        })
    })

    // Run the server on port 8080
    r.Run() // defaults to ":8080"
}
```

#### Gin Core Components
- **Engine**: The central component that manages routes and middleware
- **Context**: Carries request details and provides response methods
- **RouterGroup**: Enables route grouping for better organization
- **Middleware**: Functions that process requests before/after the main handlers

### Routing in Gin

Gin provides an intuitive API for defining routes with different HTTP methods:

```go
func main() {
    r := gin.Default()

    // Basic routes with different HTTP methods
    r.GET("/users", getUsers)
    r.POST("/users", createUser)
    r.PUT("/users/:id", updateUser)
    r.DELETE("/users/:id", deleteUser)
    
    // Route with path parameters
    r.GET("/users/:id", func(c *gin.Context) {
        id := c.Param("id")
        c.JSON(http.StatusOK, gin.H{
            "id": id,
            "message": "User details retrieved",
        })
    })
    
    // Route with query parameters
    r.GET("/search", func(c *gin.Context) {
        query := c.DefaultQuery("q", "default search")
        page := c.DefaultQuery("page", "1")
        c.JSON(http.StatusOK, gin.H{
            "query": query,
            "page": page,
        })
    })

    r.Run()
}

func getUsers(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "users": []string{"user1", "user2"},
    })
}

// Other handler functions
```

#### Route Groups

Organizing related routes into groups improves code structure:

```go
func main() {
    r := gin.Default()

    // API route group
    api := r.Group("/api")
    {
        // /api/users
        api.GET("/users", getUsers)
        
        // User-specific group
        users := api.Group("/users")
        {
            users.GET("/:id", getUserByID)
            users.PUT("/:id", updateUser)
            users.DELETE("/:id", deleteUser)
        }
        
        // Products group
        products := api.Group("/products")
        {
            products.GET("/", getProducts)
            products.POST("/", createProduct)
        }
    }

    // Admin route group with different middleware
    admin := r.Group("/admin", AuthMiddleware())
    {
        admin.GET("/analytics", getAnalytics)
        admin.GET("/users", adminGetUsers)
    }

    r.Run()
}
```

### Working with Request Data

Gin makes it easy to handle various types of request data:

```go
// User struct for binding request data
type User struct {
    ID       string `json:"id" binding:"required"`
    Username string `json:"username" binding:"required,min=4,max=20"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"required,gte=18"`
}

func createUser(c *gin.Context) {
    var user User
    
    // Bind JSON request body to user struct
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    // Process the validated user data
    // ...
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "User created successfully",
        "user": user,
    })
}

func formHandler(c *gin.Context) {
    // Parse form data
    name := c.PostForm("name")
    email := c.DefaultPostForm("email", "default@example.com")
    
    c.JSON(http.StatusOK, gin.H{
        "name": name,
        "email": email,
    })
}

func uploadHandler(c *gin.Context) {
    // Handle file upload
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    // Save the file
    dst := "uploads/" + file.Filename
    if err := c.SaveUploadedFile(file, dst); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "File uploaded successfully",
        "filename": file.Filename,
    })
}
```

### Response Handling

Gin offers various methods for sending different types of responses:

```go
func responseExamples(c *gin.Context) {
    // JSON response
    c.JSON(http.StatusOK, gin.H{
        "message": "This is a JSON response",
        "status": "success",
    })
    
    // XML response
    c.XML(http.StatusOK, gin.H{
        "message": "This is an XML response",
        "status": "success",
    })
    
    // String response
    c.String(http.StatusOK, "This is a plain text response")
    
    // HTML response (using templates)
    c.HTML(http.StatusOK, "index.html", gin.H{
        "title": "Gin HTML Template",
        "message": "Welcome to Gin!",
    })
    
    // Redirect
    c.Redirect(http.StatusMovedPermanently, "https://example.com")
    
    // File download
    c.File("path/to/file.pdf")
    
    // Custom response with headers
    c.Header("Content-Type", "application/json")
    c.Header("X-Custom-Header", "Custom Value")
    c.JSON(http.StatusOK, gin.H{"message": "Custom headers set"})
}
```

### HTML Templates in Gin

Gin integrates with Go's template system for rendering HTML:

```go
func main() {
    r := gin.Default()
    
    // Load HTML templates
    r.LoadHTMLGlob("templates/*")
    
    // Route to render template
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{
            "title": "Gin Framework",
            "content": "Welcome to Gin Web Framework!",
        })
    })
    
    // Route for a data-driven page
    r.GET("/users", func(c *gin.Context) {
        users := []gin.H{
            {"name": "Alice", "email": "alice@example.com"},
            {"name": "Bob", "email": "bob@example.com"},
            {"name": "Charlie", "email": "charlie@example.com"},
        }
        
        c.HTML(http.StatusOK, "users.html", gin.H{
            "title": "User List",
            "users": users,
        })
    })
    
    r.Run()
}
```

With corresponding template files:

```html
<!-- templates/index.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{ .title }}</title>
</head>
<body>
    <h1>{{ .title }}</h1>
    <p>{{ .content }}</p>
</body>
</html>

<!-- templates/users.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{ .title }}</title>
</head>
<body>
    <h1>{{ .title }}</h1>
    <ul>
        {{ range .users }}
        <li>{{ .name }} - {{ .email }}</li>
        {{ end }}
    </ul>
</body>
</html>
```

### Middleware in Gin

Middleware functions process requests before and after handlers, enabling cross-cutting concerns:

```go
func main() {
    // Initialize with default logger and recovery middleware
    r := gin.Default()
    
    // Global middleware - applied to all routes
    r.Use(CustomMiddleware())
    
    // Group-specific middleware
    authorized := r.Group("/", AuthMiddleware())
    {
        authorized.GET("/profile", getProfile)
    }
    
    // Route-specific middleware
    r.GET("/admin", AdminMiddleware(), adminHandler)
    
    r.Run()
}

// Custom middleware example
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Before request - preprocessing
        startTime := time.Now()
        
        // Add data to the context for handlers to use
        c.Set("example", "value")
        
        // Continue to the next middleware/handler
        c.Next()
        
        // After request - postprocessing
        latency := time.Since(startTime)
        status := c.Writer.Status()
        
        log.Printf("Request processed: %s %s %d %v",
            c.Request.Method, c.Request.URL.Path, status, latency)
    }
}

// Authentication middleware
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization required",
            })
            return
        }
        
        // Validate token (simplified)
        if token != "valid-token" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token",
            })
            return
        }
        
        // Set user info to context
        c.Set("user_id", "user123")
        
        // Continue
        c.Next()
    }
}
```

### Testing Gin Applications

Gin makes it easy to write tests for your HTTP handlers:

```go
package main

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func setupRouter() *gin.Engine {
    r := gin.Default()
    
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    
    return r
}

func TestPingRoute(t *testing.T) {
    // Set Gin to test mode
    gin.SetMode(gin.TestMode)
    
    // Create a test router
    router := setupRouter()
    
    // Create a test HTTP recorder
    w := httptest.NewRecorder()
    
    // Create a test request
    req, _ := http.NewRequest("GET", "/ping", nil)
    
    // Perform the request
    router.ServeHTTP(w, req)
    
    // Assert status code
    assert.Equal(t, http.StatusOK, w.Code)
    
    // Assert the response body
    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.Nil(t, err)
    assert.Equal(t, "pong", response["message"])
}

func TestCreateUserRoute(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    router := setupRouter()
    
    // Test JSON payload
    payload := `{"id":"123","username":"testuser","email":"test@example.com","age":25}`
    
    // Create request
    req, _ := http.NewRequest("POST", "/users", strings.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    
    // Perform request
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.Nil(t, err)
    assert.Equal(t, "User created successfully", response["message"])
}
```

### Error Handling in Gin

Proper error handling is crucial for robust web applications:

```go
func main() {
    r := gin.Default()
    
    // Custom error handling
    r.Use(ErrorMiddleware())
    
    // Routes that might trigger errors
    r.GET("/user/:id", getUserByID)
    r.POST("/user", createUser)
    
    r.Run()
}

// Application errors
type AppError struct {
    Code    int
    Message string
}

func (e *AppError) Error() string {
    return e.Message
}

// Error middleware
func ErrorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        // Check if there were any errors during handling
        if len(c.Errors) > 0 {
            for _, e := range c.Errors {
                // Check for custom error type
                if appErr, ok := e.Err.(*AppError); ok {
                    c.JSON(appErr.Code, gin.H{
                        "error": appErr.Message,
                    })
                    return
                }
            }
            
            // Default error handling
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "An unexpected error occurred",
            })
        }
    }
}

// Handler with custom error
func getUserByID(c *gin.Context) {
    id := c.Param("id")
    
    if id == "0" {
        // Add error to context
        err := &AppError{
            Code:    http.StatusNotFound,
            Message: "User not found",
        }
        c.Error(err)
        return
    }
    
    // Normal response if no error
    c.JSON(http.StatusOK, gin.H{
        "id": id,
        "name": "User Name",
    })
}
```

### Authentication and Authorization

Implementing user authentication in Gin:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
    "net/http"
    "time"
)

// Secret key for JWT
var jwtKey = []byte("my_secret_key")

// User credentials
type Credentials struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// JWT claims struct
type Claims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}

func main() {
    r := gin.Default()
    
    // Login route
    r.POST("/login", login)
    
    // Protected routes
    authorized := r.Group("/")
    authorized.Use(AuthMiddleware())
    {
        authorized.GET("/profile", getProfile)
        authorized.GET("/refresh", refreshToken)
    }
    
    // Admin routes
    admin := r.Group("/admin")
    admin.Use(AuthMiddleware(), AdminMiddleware())
    {
        admin.GET("/dashboard", adminDashboard)
    }
    
    r.Run()
}

func login(c *gin.Context) {
    var creds Credentials
    
    if err := c.ShouldBindJSON(&creds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Check credentials (simplified)
    var userRole string
    if creds.Username == "admin" && creds.Password == "admin123" {
        userRole = "admin"
    } else if creds.Username == "user" && creds.Password == "user123" {
        userRole = "user"
    } else {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // Create token expiration
    expirationTime := time.Now().Add(15 * time.Minute)
    
    // Create claims
    claims := &Claims{
        Username: creds.Username,
        Role:     userRole,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }
    
    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "token": tokenString,
        "expires_at": expirationTime,
    })
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            return
        }
        
        // Remove "Bearer " prefix if present
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }
        
        // Parse and validate token
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            return
        }
        
        // Store user info in context
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        
        c.Next()
    }
}

func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists || role != "admin" {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
            return
        }
        
        c.Next()
    }
}

func getProfile(c *gin.Context) {
    username, _ := c.Get("username")
    role, _ := c.Get("role")
    
    c.JSON(http.StatusOK, gin.H{
        "username": username,
        "role": role,
    })
}

func refreshToken(c *gin.Context) {
    // Implement token refresh logic
}

func adminDashboard(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "Welcome to the admin dashboard",
    })
}
```

### Performance Considerations

Optimizing Gin applications for production:

```go
func main() {
    // Set Gin to release mode for production
    gin.SetMode(gin.ReleaseMode)
    
    // Create a custom Engine without default middleware
    r := gin.New()
    
    // Add essential middleware only
    r.Use(gin.Recovery())
    r.Use(gin.Logger())
    
    // Add custom middleware for performance tracking
    r.Use(PerformanceMiddleware())
    
    // Configure routes
    // ...
    
    // Set up custom HTTP server with timeouts
    server := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }
    
    // Graceful shutdown
    go func() {
        quit := make(chan os.Signal, 1)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
        <-quit
        
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        if err := server.Shutdown(ctx); err != nil {
            log.Fatal("Server forced to shutdown:", err)
        }
        
        log.Println("Server exiting")
    }()
    
    // Start the server
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatal("Server startup error:", err)
    }
}

func PerformanceMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()
        
        // Process request
        c.Next()
        
        // Calculate latency
        latency := time.Since(startTime)
        
        // Log if request takes too long
        if latency > time.Second {
            log.Printf("Slow request: %s %s took %v", 
                c.Request.Method, c.Request.URL.Path, latency)
        }
    }
}
```

### Common Gin Patterns and Best Practices

1. **Proper Project Structure**
```
├── main.go           # Entry point
├── config/           # Configuration management
├── controllers/      # HTTP handlers
├── middleware/       # Custom middleware
├── models/           # Data models
├── routes/           # Route definitions
├── services/         # Business logic
├── templates/        # HTML templates
├── utils/            # Helper functions
└── tests/            # Test files
```

2. **Dependency Injection**
```go
// Service interface
type UserService interface {
    GetUser(id string) (*User, error)
    CreateUser(user *User) error
}

// Controller with dependency injection
type UserController struct {
    service UserService
}

func NewUserController(service UserService) *UserController {
    return &UserController{service: service}
}

func (uc *UserController) GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := uc.service.GetUser(id)
    
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

func SetupRoutes(r *gin.Engine, uc *UserController) {
    r.GET("/users/:id", uc.GetUser)
    // Other routes
}

func main() {
    r := gin.Default()
    
    // Initialize dependencies
    userService := NewRealUserService()
    userController := NewUserController(userService)
    
    // Setup routes
    SetupRoutes(r, userController)
    
    r.Run()
}
```

### Common Challenges and Solutions

1. **Handling CORS**
```go
func main() {
    r := gin.Default()
    
    // CORS middleware
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })
    
    // Routes
    // ...
    ge branches
    r.Run()
}
```

2. **Rate Limiting**
```go
func RateLimiter() gin.HandlerFunc {
    // Simple in-memory rate limiter
    limits := make(map[string]int)
    mutex := &sync.Mutex{}
    
    return func(c *gin.Context) {
        ip := c.ClientIP()
        
        mutex.Lock()
        if limits[ip] >= 100 {  // 100 requests per minute
            mutex.Unlock()
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
            })
            return
        }
        
        limits[ip]++
        mutex.Unlock()
        
        // Reset counters periodically (in a real app, use a timer)
        
        c.Next()
    }
}
```

### Learning Challenges

1. Create a RESTful API with Gin for a blog application (posts, comments, users)
2. Build an authentication system with JWT
3. Implement a file upload service with progress tracking
4. Create a real-time chat application using Gin and WebSockets
5. Build a microservice architecture using multiple Gin services

### Recommended Resources
- "Building Web Applications with Go and Gin" by Sam Thorogood
- Gin Framework official documentation
- "Web Development with Go" by Shiju Varghese (includes Gin sections)
- "Advanced Web Development in Go" course on Pluralsight

### Reflection Questions

1. How does Gin improve upon Go's standard library for web development?
2. What are the trade-offs between using a framework like Gin versus a microframework or the standard library?
3. How would you design a large-scale application architecture using Gin?
4. What middleware would you consider essential for a production Gin application?

**Gin Framework Mastery: Build Fast, Feature-Rich Web Applications in Go!** 🚀