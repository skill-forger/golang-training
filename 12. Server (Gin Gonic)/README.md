# Module 12: Gin Gonic

## Table of Contents

<ol>
    <li><a href="#objectives">Objectives</a></li>
    <li><a href="#overview">Overview</a></li>
    <li><a href="#introduction-to-gin-gonic">Introduction to Gin Gonic</a></li>
    <li><a href="#core-concepts">Core Concepts</a></li>
    <li><a href="#quick-start">Quick Start</a></li>
    <li><a href="#features">Features</a></li>
    <li><a href="#common-gin-patterns-and-best-practices">Common Gin Patterns and Best Practices</a></li>
    <li><a href="#practical-exercises">Practical Exercises</a></li>
</ol>

## Objectives

By the end of this module, you will be able to:

- Understand the core architecture and components of the Gin framework.
- Build robust REST APIs with routing, request handling, and validation.
- Implement middleware for cross-cutting concerns like logging and authentication.
- Structure a Gin application using best practices for scalability and maintainability.
- Write effective tests for your Gin handlers.

## Overview

Go's standard net/http library is powerful but can be verbose for building complex APIs.
Gin Gonic is a minimalistic, high-performance web framework that simplifies this process significantly.
It provides a robust set of features like routing, middleware, and rendering, allowing you to build production-ready
services quickly without sacrificing the performance that Go is known for.

## Introduction to Gin Gonic

While Go's standard library provides robust HTTP server capabilities, the Gin Gonic has emerged as one of the most
popular web frameworks in the Go ecosystem. Gin offers a more feature-rich, expressive, and performance-focused approach
to building web applications while maintaining Go's simplicity and efficiency.
For more information, refer to [Gin Gonic official documentation](https://gin-gonic.com/en/docs/).

### The Gin Landscape

Before diving into Gin's specifics, let's understand why it has become a leading choice:

#### Gin vs Standard Library

- **Performance**: Gin is built for speed with radix tree-based routing
- **Developer Experience**: More intuitive API design and middleware system
- **Feature Set**: Built-in support for JSON validation, error management, and more

#### Why Gin Matters for Modern Applications

- Simplified API development process
- Robust middleware ecosystem
- Excellent performance characteristics
- Clean code organization through groups and routes

### Syntax Comparison

#### Standard `net/http`

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.NotFound(w, r)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

#### Gin Gonic

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// gin.Default() comes with Logger and Recovery middleware.
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.Run(":8080") // listen and serve on 0.0.0.0:8080
}
```

## Core Concepts

### Engine

The Gin Engine is the core of the framework. It's the main instance that you create, and it's responsible for
registering routes, attaching middleware, and starting the server. You'll typically start with `gin.Default()` which
includes logging and panic recovery middleware.

### Context (gin.Context)

This is arguably the most important component. The `Context` is a struct that carries request-scoped data through the
middleware chain to your handler. It's how you access request details (headers, parameters, body) and how you write the
response. It's a wrapper around `http.ResponseWriter` and `*http.Request`.

### RouterGroup

This allows you to group related routes under a common path prefix and apply a shared set of middleware to all routes
within that group. It's essential for organizing your API (e.g., `/api/v1/users`, `/api/v1/products`).

### Middleware

Middleware are functions that are executed before or after the main request handler. They are arranged in a chain.
Common use cases include logging, authentication, authorization, error handling, and data transformation. A middleware
can either pass the request to the next handler in the chain using `c.Next()` or abort the request using `c.Abort()`.

## Quick Start

### Prerequisite

- Go 1.16 or above

### Installation

To install Gin package, you need to install Go and set your Go workspace first.
If you don’t have a go.mod file, create it with `go mod init gin`.

1. Download and install Gin:
    ```shell
    go get -u github.com/gin-gonic/gin
    ```
2. Import dependencies:
    ```go
    import "github.com/gin-gonic/gin"
    import "net/http"
    ```
3. Create a simple server with Gin
    ```go
   // file main.go
    package main
    
    import (
        "github.com/gin-gonic/gin"
        "net/http"
    )
    
    func main() {
        router := gin.Default()
    
        router.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
                "message": "pong",
        })
    })
    
         router.Run() // listen and serve on 0.0.0.0:8080
    }
    ```
4. Start the server
    ```shell
    go run main.go
    ```

## Features

### Routing and Handling Request

Routing is the process of mapping an incoming request's URL and HTTP method to a specific handler function.

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	// Map a GET request to the "/ping" path
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Map a POST request
	router.POST("/users", func(c *gin.Context) {
		// ... logic to create a user
		c.JSON(http.StatusCreated, gin.H{"status": "user created"})
	})

	router.Run(":8080")
}

```

Gin provides an intuitive API for defining routes with different HTTP methods,
supports all standard HTTP methods: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, and `OPTIONS`.

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
			"id":      id,
			"message": "User details retrieved",
		})
	})

	// Route with query parameters
	r.GET("/search", func(c *gin.Context) {
		query := c.DefaultQuery("q", "default search")
		page := c.DefaultQuery("page", "1")
		c.JSON(http.StatusOK, gin.H{
			"query": query,
			"page":  page,
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

### Working with Route Groups

Organizing related routes into groups improves code structure.

```go
package main

import (
	"github.com/gin-gonic/gin"
)

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

Router Groups are perfect for versioning your API and applying shared logic.

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	// Group for API version 1
	v1 := router.Group("/api/v1")
	{ // Use braces for visual separation
		v1.GET("/users", GetUsers)       // Handler for /api/v1/users
		v1.POST("/users", CreateUser)    // Handler for /api/v1/users
		v1.GET("/products", GetProducts) // Handler for /api/v1/products
	}

	// Group for API version 2
	v2 := router.Group("/api/v2")
	{
		v2.GET("/users", GetUsersV2)
	}

	router.Run(":8080")
}

// Dummy handler functions
func GetUsers(c *gin.Context)    { c.JSON(http.StatusOK, "v1 users") }
func CreateUser(c *gin.Context)  { c.JSON(http.StatusCreated, "v1 user created") }
func GetProducts(c *gin.Context) { c.JSON(http.StatusOK, "v1 products") }
func GetUsersV2(c *gin.Context)  { c.JSON(http.StatusOK, "v2 users") }

```

### Binding Payload and parsing data

Gin makes it trivial to extract data from a request.

#### Path, Form data and Query Parameters

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	// Path Parameter (e.g., /users/123)
	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, "User ID is %s", id)
	})

	// Query Parameter (e.g., /search?query=golang)
	router.GET("/search", func(c *gin.Context) {
		// Use Query for a single value, or DefaultQuery for a fallback
		query := c.DefaultQuery("query", "guest")
		c.String(http.StatusOK, "Search query is '%s'", query)
	})

	// Form data (e.g., /form-submit)
	router.POST("/form-submit", func(c *gin.Context) {
		// Parse form data
		name := c.PostForm("name")
		email := c.DefaultPostForm("email", "default@example.com")

		c.JSON(http.StatusOK, gin.H{
			"name":  name,
			"email": email,
		})
	})

	router.Run(":8080")
}

```

#### Request Body

Binding automatically parses the request body (e.g., JSON) into a Go struct. This is extremely powerful for validation.

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
		"user":    user,
	})
}

```

#### Files and Assets

Gin allows API to handle files and assets uploading.

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
		"message":  "File uploaded successfully",
		"filename": file.Filename,
	})
}

```

Gin can render HTML templates and serve static files like CSS, JS, and images.

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	// Serve static files from the "./static" directory
	router.Static("/static", "./static")

	// Load HTML templates from the "./templates" directory
	router.LoadHTMLGlob("templates/*")

	router.GET("/index", func(c *gin.Context) {
		// Render the "index.tmpl" template
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.Run(":8080")
}

```

With corresponding template files.

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

### Handling Response

Gin offers various methods for sending different types of responses:

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func responseExamples(c *gin.Context) {
	// JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "This is a JSON response",
		"status":  "success",
	})

	// XML response
	c.XML(http.StatusOK, gin.H{
		"message": "This is an XML response",
		"status":  "success",
	})

	// String response
	c.String(http.StatusOK, "This is a plain text response")

	// HTML response (using templates)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":   "Gin HTML Template",
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

### Applying Middlewares

Middleware allows you to inject logic into the request-processing pipeline.
Let's create a simple authentication middleware.

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware checks for a specific API token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-API-Token")

		// In a real app, you'd validate this token properly.
		if token != "super-secret-token" {
			// Abort the request and send an error response.
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API token required"})
			return
		}

		// Token is valid, pass the request to the next handler in the chain.
		c.Next()
	}
}

func main() {
	router := gin.Default()

	// Secure group using the middleware
	secured := router.Group("/secure")
	secured.Use(AuthMiddleware()) // Apply middleware to the group
	{
		secured.GET("/profile", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "this is a secure endpoint"})
		})
	}

	router.Run(":8080")
}

```

### Chaining Middleware

In Gin, middleware functions are simply handler functions with a special purpose. They have the signature `func(c *
gin.Context)`. The key to chaining is the `c.Next()` method. When a middleware function calls `c.Next()`, it pauses its
own execution and passes control to the next middleware or the final handler in the chain. After the next function
completes, control returns to the original middleware, which can then perform actions (e.g., logging a request's
duration).

If a middleware does not call `c.Next()`, the chain is effectively "short-circuited." The request will not proceed to
subsequent middleware or the final handler, and the middleware is responsible for sending a response to the client. This
is how authentication middleware, for example, can stop a request from reaching a protected route if the user is not
authorized.

**Example 1**: Basic Chaining with `c.Next()`
In this example, when a GET request hits `/`, Gin first executes our `LoggerMiddleware`. This middleware records the
start
time and then calls `c.Next()`. This passes control to the final handler function, which sends a JSON response. After
the
handler is finished, control returns to the `LoggerMiddleware`, which calculates the request's latency and logs it.

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Logging middleware
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// Process the next function in the chain
		c.Next()

		// Code here runs after the handler and subsequent middleware have completed
		latency := time.Since(start)
		log.Printf("Request took %v | %s | %s", latency, c.Request.Method, c.Request.URL.Path)
	}
}

func main() {
	r := gin.Default() // gin.Default() comes with Logger and Recovery middleware built-in

	// Apply our custom logging middleware to the "/" route
	r.GET("/", LoggerMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from Gin!"})
	})

	r.Run(":8080")
}
```

**Example 2**: Chaining with Authentication (Short-Circuiting)
This example shows how a middleware can abort a request if a condition isn't met,
preventing the final handler from being executed. In the code snippet below:

- The `/public` route is open to everyone.
- The `/protected` route group uses `AuthMiddleware()`. Any request to `/protected/admin` will first go through this
  middleware.
- The `AuthMiddleware` checks for a specific authorization token. If the token is incorrect, it
  calls `c.AbortWithStatusJSON()`. This both sets the response and stops the request from going any further in the
  chain. The final handler `protected.GET("/admin", ...)` is never reached.
- If the token is valid, `c.Next()` is called, and the request proceeds to the handler, which then returns a successful
  response.

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for an 'Authorization' header
		token := c.GetHeader("Authorization")
		if token != "Bearer mysecrettoken" {
			// If not authorized, abort the request and send a 401 response
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return // Don't call c.Next()
		}
		// If authorized, continue to the next handler
		c.Next()
	}
}

func main() {
	r := gin.Default()

	// Public route - no middleware
	r.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is a public route"})
	})

	// Protected route group
	// Apply the AuthMiddleware to all routes in this group
	protected := r.Group("/protected")
	protected.Use(AuthMiddleware()) // .Use() is for applying middleware to a group

	{
		protected.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Welcome, authorized user!"})
		})
	}

	r.Run(":8080")
}

```

**Example 3**: Chaining Multiple Middlewares

In the code snippet below:

1. Request Flow: A request to `/api/protected/resource` first enters the `protected` group.
2. `RequestLogger`: This middleware runs first. It logs the start of the request and calls `c.Next()`, passing control
   to
   the next middleware in the chain.
3. `RateLimiter`: This middleware runs second. It checks if the client has exceeded the request limit. If the limit is
   exceeded, it calls `c.AbortWithStatusJSON()` and returns, effectively stopping the chain. If not, it calls `c.Next()`
   to
   continue.
4. `AuthMiddleware`: This middleware runs third. It checks for a valid `Authorization` header. If the token is
   incorrect, it
   calls `c.AbortWithStatusJSON()` and returns, stopping the chain. If the token is valid, it calls `c.Next()`.
5. Final Handler: Only if all three middlewares successfully call `c.Next()` will the request reach the final handler
   for
   `/resource`. The handler then sends its JSON response. After the response is sent, execution returns up the chain (in
   reverse order), allowing the `RequestLogger` to complete its timing and logging.

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger middleware logs the request method and path.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // Pass control to the next middleware or handler

		// This part runs after the handler has completed
		latency := time.Since(start)
		log.Printf("[RequestLogger] %s %s took %v", c.Request.Method, c.Request.URL.Path, latency)
	}
}

// RateLimiter middleware limits requests to a specific endpoint.
func RateLimiter() gin.HandlerFunc {
	// A simple rate limit map to simulate limiting by client IP
	requests := make(map[string]int)
	lastReset := time.Now()

	return func(c *gin.Context) {
		// Reset the counter every minute
		if time.Since(lastReset) > 1*time.Minute {
			requests = make(map[string]int)
			lastReset = time.Now()
		}

		clientIP := c.ClientIP()
		requests[clientIP]++
		if requests[clientIP] > 5 { // Allow up to 5 requests per minute
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			log.Printf("[RateLimiter] Rate limit exceeded for IP: %s", clientIP)
			return
		}
		c.Next()
	}
}

// AuthMiddleware checks for a valid authorization token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer mysecrettoken" {
			// Abort the chain if the token is invalid
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			log.Printf("[AuthMiddleware] Unauthorized request with token: %s", token)
			return
		}
		log.Printf("[AuthMiddleware] User authorized")
		c.Next()
	}
}

func main() {
	router := gin.Default()

	// Apply multiple middlewares to a route group using .Use()
	protected := router.Group("/api/protected")
	protected.Use(RequestLogger(), RateLimiter(), AuthMiddleware())
	{
		// This handler will only be reached if all three middlewares pass
		protected.GET("/resource", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Access granted to the protected resource!"})
		})
	}

	// This is a public route that bypasses all the middlewares
	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is a public endpoint."})
	})

	router.Run(":8080")
}
```

### Handling Error

A centralized error-handling middleware is a robust pattern. It catches errors returned from your handlers and formats a
consistent error response.

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Custom error struct
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

// Error handling middleware
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process request

		// After handler, check for errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appErr *AppError
			if errors.As(err, &appErr) {
				// This is our custom error
				c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			} else {
				// This is an unexpected internal error
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}
	}
}

// Handler that can return an error
func getUser(c *gin.Context) {
	// ... logic to get a user ...
	// If user not found:
	err := &AppError{Code: http.StatusNotFound, Message: "User not found"}
	c.Error(err) // Push error to context
	c.Abort()    // Stop processing
}

```

### Testing

#### Set up the Test Environment

Before writing tests, it's good practice to create a test setup function. This function will instantiate a Gin router
and configure its routes, ensuring that each test starts with a clean slate. You'll typically place this code in a file
ending with `_test.go` (e.g., `main_test.go`).

A common pattern is to create a function that returns the router instance:

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// setupRouter initializes the Gin router with all its routes.
func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return router
}
```

**Tip**: Gin has a gin.TestMode that you can set to disable its debug output and color logging during tests.
This helps keep your test output clean. A TestMain function is the ideal place for this.

```go
package main

import (
	"github.com/gin-gonic/gin"
	"testing"
)

// in main_test.go
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

```

#### Testing a GET Endpoint

Testing a `GET` request is the most straightforward scenario. You don't need to send a request body, only the path.

Scenario: Test a simple `/ping` endpoint that returns a JSON response.

```go
package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestPingRoute tests the GET /ping endpoint
func TestPingRoute(t *testing.T) {
	// 1. Get the router instance
	router := setupRouter()

	// 2. Create a mock response recorder
	w := httptest.NewRecorder()

	// 3. Create a mock HTTP request
	req, _ := http.NewRequest("GET", "/ping", nil)

	// 4. Serve the HTTP request to the router
	router.ServeHTTP(w, req)

	// 5. Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

```

#### Testing a POST Endpoint with JSON Body

Testing a `POST` request requires sending a request body and setting the correct content type header.

Scenario: Test a `/user` endpoint that binds a JSON body to a struct and returns the created user.

First, define the handler logic:

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// User represents a user struct.
type User struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email"`
}

// postUserHandler is a handler for creating a new user.
func postUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/user", postUserHandler)
	return router
}
```

Testing logic:

```go
package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// in main_test.go
func TestPostUser(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()

	// Define the request body as a JSON string
	jsonBody := `{"name":"John Doe", "email":"john.doe@example.com"}`
	req, _ := http.NewRequest("POST", "/user", strings.NewReader(jsonBody))

	// Set the Content-Type header to simulate a real request
	req.Header.Set("Content-Type", "application/json")

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, jsonBody, w.Body.String())
}

```

#### Testing with URL and Query Parameters

**Scenario**: Test an endpoint that uses both a URL parameter and a query parameter.

```go
package main

// in main.go
func getProductHandler(c *gin.Context) {
	id := c.Param("id")
	sort := c.DefaultQuery("sort", "asc")
	c.JSON(http.StatusOK, gin.H{
		"product_id": id,
		"sort_by":    sort,
	})
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/products/:id", getProductHandler)
	return router
}
```

Now, the test:

```go
package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// in main_test.go
func TestGetProductWithParams(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()

	// Test with both URL and query parameters
	req, _ := http.NewRequest("GET", "/products/123?sort=desc", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"product_id":"123", "sort_by":"desc"}`, w.Body.String())

	// Test with URL parameter only (to check the default query)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/products/456", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"product_id":"456", "sort_by":"asc"}`, w.Body.String())
}

```

## Common Gin Patterns and Best Practices

#### Project Structure

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

### Dependency Injection

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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

## Common Challenges and Solutions

### Handling CORS

```go
package main

import "github.com/gin-gonic/gin"

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
	ge
	branches
	r.Run()
}

```

### Rate Limiting

```go
package main

import "github.com/gin-gonic/gin"

func RateLimiter() gin.HandlerFunc {
	// Simple in-memory rate limiter
	limits := make(map[string]int)
	mutex := &sync.Mutex{}

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mutex.Lock()
		if limits[ip] >= 100 { // 100 requests per minute
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

## Practical Exercises

### Exercise 1: Basic Gin API

Create a simple RESTful API using Gin framework to manage a todo list:

### Exercise 2: Gin Middleware and Authentication

Create a Gin application with custom middleware for logging and simple API key authentication:

### Exercise 3: File Upload with Gin

Create a Gin application that handles file uploads with progress monitoring: