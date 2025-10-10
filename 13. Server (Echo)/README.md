# Module 12: Echo

## Table of Contents

<ol>
    <li><a href="#objectives">Objectives</a></li>
    <li><a href="#overview">Overview</a></li>
    <li><a href="#introduction-to-echo">Introduction to Echo</a></li>
    <li><a href="#core-concepts">Core Concepts</a></li>
    <li><a href="#quick-start">Quick Start</a></li>
    <li><a href="#features">Features</a></li>
    <li><a href="#common-echo-patterns-and-best-practices">Common Echo Patterns and Best Practices</a></li>
    <li><a href="#practical-exercises">Practical Exercises</a></li>
</ol>

## Objectives

By the end of this module, you will be able to:

- Understand Echo fundamentals and set up a basic project.
- Build RESTful APIs with routing and CRUD operations.
- Use built-in middleware and create custom middleware.
- Optimize performance and deploy applications to production.
- Write tests and apply practices in Echo projects.

## Overview

Echo is a modern and lightweight web framework for Go that makes it simple to build fast, scalable, and secure web
applications.
With its clean design, powerful routing, and rich middleware support, Echo helps developers quickly create RESTful APIs
and backend
services while keeping code maintainable and efficient. This module guides you step by step from the basics to advanced topics, 
helping you master Echo through practical examples.

## Introduction to Echo

The Echo project offers an array of features that empower developers to build robust web applications.
Its fast and lightweight nature ensures optimal performance, while the flexible routing system and middleware support
streamline development processes.
Developers can leverage the context-based request handling, powerful template rendering, and validation capabilities to
create dynamic and secure web applications.
Additionally, the extensibility of Echo allows developers to customize and enhance the framework to suit their specific
needs.
For more information, refer to [Echo official documentation](https://echo.labstack.com/docs).

### The Echo Landscape

Before diving into Echo’s specifics, let’s understand why it has become one of the most popular Go frameworks:

#### Echo vs Standard Library

- **Performance**: Echo delivers lightning-fast routing with minimal overhead.
- **Developer Experience**: Clean and simple API design that reduces boilerplate code.
- **Feature Set**: Rich middleware support, built-in request binding, and flexible error handling.

#### Why Echo Matters for Modern Applications

- Speeds up building RESTful APIs and microservices.
- Provides a strong middleware ecosystem with easy customization.
- Designed for scalability and production-ready performance.
- Encourages maintainable code structure with groups and modular routes.

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

#### Echo

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()

	router.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	router.Start(":8080") // listen and serve on 0.0.0.0:8080
}
```

## Core Concepts

### Engine

The Echo Engine is the core of the framework. It’s the main instance you create with `echo.New()`. 
The engine is responsible for registering routes, attaching middleware, and starting the HTTP server. 
You can customize it by adding middleware, defining handlers, and configuring server settings.

### Echo Context (`echo.Context`)

The `Context` is the most important component in Echo. It encapsulates the HTTP request and response, 
giving you convenient methods to read parameters, headers, query strings, and request bodies. It also 
provides functions to send responses in various formats (JSON, HTML, text, file, etc.), making request–response 
handling simple and efficient.

### Router Groups

Echo provides Group functionality to organize routes under a shared path prefix and middleware. For example, you can
group routes by version (e.g., `/api/v1/users`, `/api/v1/products`) and apply authentication middleware only once at 
the group level. This keeps code modular and easier to maintain.

### Middleware

Middleware in Echo are functions executed before or after handlers. They can be applied globally, per route, or per
group. Common use cases include logging, error handling, authentication, authorization, rate limiting, and CORS. Middleware 
can control request flow by calling `next(c)` to continue or aborting with an error response.

## Quick Start

### Prerequisite

- Go 1.13 or higher. Go 1.12 has limited support and some middlewares will not be available

### Installation

To install Echo package, you need to install Go and set your Go workspace first.
If you don’t have a go.mod file, create it with:

```
$ mkdir myapp && cd myapp
$ go mod init myapp
```

1. Download and install Gin:

```shell
go get -u github.com/labstack/echo/v4
```

2. Create `server.go`
3. Import dependencies:

```go
import (
"net/http"

"github.com/labstack/echo/v4"
)
```

4. Create a simple server with Echo

```go
// file main.go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()

	router.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "pong",
		})
	})

	router.Start(":8080") // listen and serve on 0.0.0.0:8080
}
```

4. Start the server

```shell
go run server.go
```

## Features

### Routing and Handling Request

Routing is the process of mapping an incoming request's URL and HTTP method to a specific handler function.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()

	// Map a GET request to the "/ping" path
	router.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "pong",
		})
	})

	// Map a PUT request to the "/users/:id" path
	router.PUT("/users/:id", func(c echo.Context) error {
		// ... logic to update a user by id
		return c.JSON(http.StatusOK, map[string]any{
			"status": "user updated",
		})
	})

	router.Start(":8080") // listen and serve on 0.0.0.0:8080
}
```

Echo provides an intuitive API for defining routes with different HTTP methods,
supports all standard HTTP methods: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, and `OPTIONS`.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()

	// Basic routes with different HTTP methods
	router.GET("/users", getUsers)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	// Route with path parameters
	router.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.JSON(http.StatusOK, map[string]any{
			"id":      id,
			"message": "Fetched user details successfully",
		})
	})

	// Route with query parameters
	router.GET("/search", func(c echo.Context) error {
		query := c.QueryParam("q")
		if query == "" {
			query = "default search"
		}
		page := c.QueryParam("page")
		if page == "" {
			page = "1"
		}
		return c.JSON(http.StatusOK, map[string]any{
			"query": query,
			"page":  page,
		})
	})

	// Start server
	router.Start(":8080")
}

func getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"users": []string{"alice", "bob"},
	})
}

// Other handler functions
```

### Working with Route Groups

Organizing related routes into groups improves code structure.

```go
package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()

	// API route group
	api := router.Group("/api")
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
	admin := router.Group("/admin", AuthMiddleware())
	{
		admin.GET("/analytics", getAnalytics)
		admin.GET("/users", adminGetUsers)
	}

	router.Start(":8080")
}

```

Router Groups are perfect for versioning your API and applying shared logic.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()

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

	router.Start(":8080")
}

// Dummy handler functions
func GetUsers(c echo.Context) error    { return c.JSON(http.StatusOK, "v1 users") }
func CreateUser(c echo.Context) error  { return c.JSON(http.StatusCreated, "v1 user created") }
func GetProducts(c echo.Context) error { return c.JSON(http.StatusOK, "v1 products") }
func GetUsersV2(c echo.Context) error  { return c.JSON(http.StatusOK, "v2 users") }

```

### Binding Payload and parsing data

Echo makes it trivial to extract data from a request.

#### Path, Form data and Query Parameters

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()

	// Path Parameter (e.g., /users/123)
	router.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")

		return c.String(http.StatusOK, fmt.Sprintf("User ID is %s", id))
	})

	// Query Parameter (e.g., /search?query=golang)
	router.GET("/search", func(c echo.Context) error {
		query := c.QueryParam("query")

		return c.String(http.StatusOK, fmt.Sprintf("Search query is '%s'", query))
	})

	// Form data (e.g., /form-submit)
	router.POST("/form-submit", func(c echo.Context) error {
		// Parse form data
		name := c.FormValue("name")
		email := c.FormValue("email")

		return c.JSON(http.StatusOK, map[string]any{
			"name":  name,
			"email": email,
		})
	})

	router.Start(":8080")
}
```

#### Request Body

Binding automatically parses the request body (e.g., JSON) into a Go struct. This is extremely powerful for validation.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// User struct for binding request data
type User struct {
	ID       string `json:"id" validate:"required"`
	Age      int    `json:"age" validate:"required,gte=18"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=4,max=20"`
}

func createUser(c echo.Context) error {
	var user User

	// Bind JSON request body to user struct
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
	}

	// Process the validated user data
	// ...

	return c.JSON(http.StatusCreated, map[string]any{
		"message": "User created successfully",
		"user":    user,
	})
}
```

#### Files and Assets

Echo allows API to handle files and assets uploading.

```go
package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func uploadHandler(c echo.Context) error {
	// Lấy file từ form-data (key: file)
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
	}

	// Mở file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}
	defer src.Close()

	// Tạo file đích
	dstPath := "uploads/" + file.Filename
	dst, err := os.Create(dstPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}
	defer dst.Close()

	// Copy dữ liệu
	if _, err := dst.ReadFrom(src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	// Trả về response
	return c.JSON(http.StatusOK, map[string]any{
		"message":  "File uploaded successfully",
		"filename": file.Filename,
	})
}
```

Echo can render HTML templates and serve static files like CSS, JS, and images.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func upload(c echo.Context) error {
	//Upload file logic
	//...
	return c.HTML(http.StatusOK, fmt.Sprintf("<p>Uploaded files successfully.</p>"))
}

func main() {
	router := echo.New()

	router.POST("/upload", upload)

	router.Start(":8080")
}
```

### Handling Response

Echo offers various methods for sending different types of responses:

```go
package main

import (
	"github.com/labstack/echo/v4"

	"net/http"
)

func responseExamples(c echo.Context) {
	// JSON response
	c.JSON(http.StatusOK, map[string]any{
		"message": "This is a JSON response",
		"status":  "success",
	})

	// XML response
	c.XML(http.StatusOK, map[string]any{
		"message": "This is an XML response",
		"status":  "success",
	})

	// String response
	c.String(http.StatusOK, "This is a plain text response")

	// HTML response (using templates)
	c.HTML(http.StatusOK, "index.html")

	// Redirect
	c.Redirect(http.StatusMovedPermanently, "https://example.com")

	// File download
	c.File("path/to/file.pdf")

	// Custom response with headers
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("X-Custom-Header", "Custom Value")
	c.JSON(http.StatusOK, map[string]any{"message": "Custom headers set"})
}
```

### Applying Middlewares

Middleware allows you to inject logic into the request-processing pipeline.
Let's create a simple authentication middleware.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware checks for a specific API token.
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("X-API-Token")

		// In a real app, you'd validate this token properly
		if token != "super-secret-token" {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error": "API token required",
			})
		}

		// // Token is valid, pass the request to the next handler in the chain.
		return next(c)
	}
}

func main() {
	router := echo.New()

	// Secure group using the middleware
	secured := router.Group("/secure", AuthMiddleware)
	secured.GET("/profile", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "this is a secure endpoint",
		})
	})

	router.Start(":8080")
}
```

### Chaining Middleware

In Echo, middleware functions are special handler functions that can process logic 
before and/or after the main route handler executes. They have the signature:
```go
func(next echo.HandlerFunc) echo.HandlerFunc
```
A middleware receives a next handler and returns another handler.
To pass control to the next middleware or final handler, it must call:
```go
return next(c)
```
If a middleware does not call `next(c)`, the chain is effectively short-circuited — the 
request will not continue to subsequent middleware or the final route handler. In that 
case, the middleware itself is responsible for sending the response to the client.

**Example 1**: Basic Chaining with `c.Next()`
In this example, when a GET request hits `/`, Echo first executes our `LoggerMiddleware`. This 
middleware records the start time, then calls `next(c)` — passing control to the final handler 
function, which sends a JSON response. After the handler completes, control returns to the 
middleware, which calculates the request’s latency and logs it.

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Logging middleware
func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process the next function in the chain
			err := next(c)

			// Code here runs after the handler and subsequent middleware have completed
			latency := time.Since(start)
			log.Printf("Request took %v | %s | %s", latency, c.Request().Method, c.Request().URL.Path)

			return err
		}
	}
}

func main() {
	router := echo.New()

	// Apply our custom logging middleware to the "/" route
	router.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "Hello from Echo!",
		})
	}, LoggerMiddleware())

	router.Logger.Fatal(e.Start(":8080"))
}
```

**Example 2**: Chaining with Authentication (Short-Circuiting)
This example shows how an Echo middleware can abort a request if a condition isn’t 
met — preventing the final handler from executing.

- The `/public` route is open to everyone.
- The `/protected` route group uses `AuthMiddleware()`. Any request to `/protected/admin` will first 
- go through this middleware.
- The `AuthMiddleware` checks for a specific authorization token.
- If the token is invalid, it returns an error response immediately, without calling `next(c)` — effectively short-circuiting the chain.
- If the token is valid, it calls `next(c)`, allowing the request to proceed to the handler, which returns a successful response.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Authentication middleware
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check for an 'Authorization' header
			token := c.Request().Header.Get("Authorization")

			if token != "Bearer mysecrettoken" {
				// If not authorized, return 401 and stop further execution
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"error": "Unauthorized",
				})
			}

			// If authorized, continue to the next handler
			return next(c)
		}
	}
}

func main() {
	router := echo.New()

	// Public route - no middleware
	router.GET("/public", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "This is a public route",
		})
	})

	// Protected route group
	protected := router.Group("/protected", AuthMiddleware())

	protected.GET("/admin", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "Welcome, authorized user!",
		})
	})

	router.Logger.Fatal(router.Start(":8080"))
}
```

**Example 3**: Chaining Multiple Middlewares

In the code snippet below:

1. Request Flow: A request to `/api/protected/resource` first enters the `protected` group.
2. `RequestLogger`: Runs first. Logs request start time, then calls `next(c)` to pass control to the next middleware.
3. `RateLimiter`: Runs second. If the client exceeds the limit, it returns a `429 Too Many Requests` response and stops the chain. Otherwise, it calls `next(c)`.
4. `AuthMiddleware`: Runs third. If the `Authorization` header is invalid, it returns `401 Unauthorized` and stops the chain. If valid, it calls `next(c)`.
5. Final Handler: Executes only if all previous middlewares successfully call `next(c)`. After the response is sent, control returns up the chain, allowing RequestLogger to complete its timing and logging.

```go
package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// RequestLogger middleware logs the request method and path.
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c) // Pass control to the next middleware or handler

			// This part runs after the handler has completed
			latency := time.Since(start)
			log.Printf("[RequestLogger] %s %s took %v", c.Request().Method, c.Request().URL.Path, latency)

			return err
		}
	}
}

// RateLimiter middleware limits requests to a specific endpoint.
func RateLimiter() echo.MiddlewareFunc {
	var mu sync.Mutex
	requests := make(map[string]int)
	lastReset := time.Now()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			mu.Lock()
			defer mu.Unlock()

			// Reset the counter every minute
			if time.Since(lastReset) > time.Minute {
				requests = make(map[string]int)
				lastReset = time.Now()
			}

			clientIP := c.RealIP()
			requests[clientIP]++

			if requests[clientIP] > 5 { // Allow up to 5 requests per minute
				log.Printf("[RateLimiter] Rate limit exceeded for IP: %s", clientIP)
				return c.JSON(http.StatusTooManyRequests, map[string]any{
					"error": "Too many requests",
				})
			}

			// Continue to the next middleware
			return next(c)
		}
	}
}

// AuthMiddleware checks for a valid authorization token.
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token != "Bearer mysecrettoken" {
				log.Printf("[AuthMiddleware] Unauthorized request with token: %s", token)
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"error": "Unauthorized",
				})
			}

			log.Printf("[AuthMiddleware] User authorized")
			return next(c)
		}
	}
}

func main() {
	router := echo.New()

	// Protected route group with multiple middlewares
	protected := router.Group("/api/protected",
		RequestLogger(),
		RateLimiter(),
		AuthMiddleware(),
	)

	// Handler only runs if all middlewares succeed
	protected.GET("/resource", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "Access granted to the protected resource!",
		})
	})

	// Public route bypassing all middlewares
	router.GET("/public", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "This is a public endpoint.",
		})
	})

	router.Logger.Fatal(router.Start(":8080"))
}
```

### Handling Error
A centralized error-handling middleware is a robust pattern. It catches errors returned from your handlers and formats a
consistent error response.

```go
package main

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
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
func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Execute the next handler in the chain
		err := next(c)
		if err == nil {
			return nil
		}

		// Check if it's our custom AppError
		var appErr *AppError
		if errors.As(err, &appErr) {
			return c.JSON(appErr.Code, map[string]any{
				"error": appErr.Message,
			})
		}

		// Unexpected internal error
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "Internal Server Error",
		})
	}
}

// Handler that can return an error
func getUser(c echo.Context) error {
	// ... logic to get a user ...
	// Simulate a "not found" situation
	return &AppError{Code: http.StatusNotFound, Message: "User not found"}
}

func main() {
	router := echo.New()

	// Apply the error handling middleware globally
	router.Use(ErrorHandler)

	// Route that can trigger an error
	router.GET("/user/:id", getUser)

	router.Logger.Fatal(e.Start(":8080"))
}
```

### Testing

#### Set up the Test Environment

Before writing tests, it’s good practice to create a reusable setup function.
This function initializes a clean Echo instance, configures routes, and returns it for each test case.
You’ll typically place this code in a file named `main_test.go` or `router_test.go`.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// setupRouter initializes the Echo router with all its routes.
func setupRouter() *echo.Echo {
	router := echo.New()

	router.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "pong",
		})
	})

	return router
}

```

#### Testing a GET Endpoint
Testing a `GET` request is the most straightforward scenario. You don't need to send a request body, only the path.

Scenario: Test a simple `/ping` endpoint that returns a JSON response.

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestPingRoute tests the GET /ping endpoint
func TestPingRoute(t *testing.T) {
	// 1. Get the Echo router instance
	router := setupRouter()

	// 2. Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	// 3. Create a new Echo context
	c := router.NewContext(req, rec)

	// 4. Serve the request to the router
	// Note: In Echo, use e.ServeHTTP(rec, req) for end-to-end testing
	router.ServeHTTP(rec, req)

	// 5. Assert the response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"pong"}`, rec.Body.String())
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

	"github.com/labstack/echo/v4"
)

// User represents a user struct.
type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email"`
}

// postUserHandler handles creating a new user.
func postUserHandler(c echo.Context) error {
	var user User

	// Bind JSON body to struct
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
	}
	
	// Return created user as JSON
	return c.JSON(http.StatusCreated, user)
}

// setupRouter initializes the Echo router with routes
func setupRouter() *echo.Echo {
	router := echo.New()
	router.POST("/user", postUserHandler)
	return router
}
```

Testing logic:

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// in main_test.go
func TestPostUser(t *testing.T) {
	router := setupRouter()

	// Define the request body as a JSON string
	jsonBody := `{"name":"John Doe", "email":"john.doe@example.com"}`

	// Create a new HTTP POST request with JSON body
	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rec := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rec, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, jsonBody, rec.Body.String())
}
```

#### Testing with URL and Query Parameters
**Scenario**: Test an endpoint that uses both a URL parameter and a query parameter.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// getProductHandler handles a request with a URL param and query param
func getProductHandler(c echo.Context) error {
	id := c.Param("id")
	sort := c.QueryParam("sort")
	if sort == "" {
		sort = "asc"
	}

	return c.JSON(http.StatusOK, map[string]any{
		"product_id": id,
		"sort_by":    sort,
	})
}

// setupRouter initializes the Echo router with routes
func setupRouter() *echo.Echo {
	router := echo.New()
	router.GET("/products/:id", getProductHandler)
	return router
}
```

Now, the test:

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// in main_test.go
func TestGetProductWithParams(t *testing.T) {
	router := setupRouter()

	// Test with both URL and query parameters
	req := httptest.NewRequest(http.MethodGet, "/products/123?sort=desc", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"product_id":"123","sort_by":"desc"}`, rec.Body.String())

	// Test with URL parameter only (no query -> default to "asc")
	req = httptest.NewRequest(http.MethodGet, "/products/456", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"product_id":"456","sort_by":"asc"}`, rec.Body.String())
}
```

## Common Echo Patterns and Best Practices

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
	"net/http"

	"github.com/labstack/echo/v4"
)

// User represents a user model
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Service interface
type UserService interface {
	GetUser(id string) (*User, error)
	CreateUser(user *User) error
}

// Controller with dependency injection
type UserController struct {
	service UserService
}

// Constructor
func NewUserController(service UserService) *UserController {
	return &UserController{service: service}
}

// Handler: GET /users/:id
func (uc *UserController) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := uc.service.GetUser(id)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]any{
			"error": "User not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// Setup routes
func SetupRoutes(e *echo.Echo, uc *UserController) {
	e.GET("/users/:id", uc.GetUser)
	// Other routes (e.g., e.POST("/users", uc.CreateUser))
}

// Example of a real implementation of UserService
type RealUserService struct{}

func NewRealUserService() *RealUserService {
	return &RealUserService{}
}

func (s *RealUserService) GetUser(id string) (*User, error) {
	// Simulate data fetch
	if id != "1" {
		return nil, echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return &User{ID: "1", Name: "John Doe", Email: "john@example.com"}, nil
}

func (s *RealUserService) CreateUser(user *User) error {
	// Simulate creating a user
	return nil
}

func main() {
	router := echo.New()

	// Initialize dependencies
	userService := NewRealUserService()
	userController := NewUserController(userService)

	// Setup routes
	SetupRoutes(router, userController)

	router.Logger.Fatal(router.Start(":8080"))
}
```

## Common Challenges and Solutions

### Handling CORS

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Custom CORS middleware
func CORSMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Set CORS headers
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
		c.Response().Header().Set(echo.HeaderAccessControlAllowMethods, "GET, POST, PUT, DELETE, OPTIONS")
		c.Response().Header().Set(echo.HeaderAccessControlAllowHeaders, "Content-Type, Authorization")

		// Handle preflight requests
		if c.Request().Method == http.MethodOptions {
			return c.NoContent(http.StatusNoContent)
		}

		// Continue to next handler
		return next(c)
	}
}

func main() {
	router := echo.New()

	// Apply CORS middleware globally
	router.Use(CORSMiddleware)

	// Example routes
	router.GET("/branches", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "List of branches",
		})
	})

	router.POST("/branches", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, map[string]any{
			"message": "Branch created",
		})
	})

	// Start the server
	router.Logger.Fatal(router.Start(":8080"))
}

```

### Rate Limiting

```go
package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// RateLimiter middleware for Echo
func RateLimiter() echo.MiddlewareFunc {
	limits := make(map[string]int)
	mutex := &sync.Mutex{}
	lastReset := time.Now()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()

			mutex.Lock()
			defer mutex.Unlock()

			// Reset counters every minute
			if time.Since(lastReset) > time.Minute {
				limits = make(map[string]int)
				lastReset = time.Now()
			}

			// Check the request count for this IP
			if limits[ip] >= 100 { // 100 requests per minute
				return c.JSON(http.StatusTooManyRequests, map[string]any{
					"error": "Rate limit exceeded",
				})
			}

			// Increment count
			limits[ip]++

			// Continue to the next middleware/handler
			return next(c)
		}
	}
}

func main() {
	router := echo.New()

	// Apply rate limiter globally
	router.Use(RateLimiter())

	router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	router.Logger.Fatal(router.Start(":8080"))
}

```

## Practical Exercises

### Exercise 1: Basic Echo API

Create a simple RESTful API using Echo framework to manage a todo list:

### Exercise 2: CRUD “User Management” API

Practice CRUD operations (Create, Read, Update, Delete):

### Exercise 3: Middleware & Validation

Create a middleware that checks for an X-API-KEY header. Validate request data before creating a user. Apply middleware to a specific route group.