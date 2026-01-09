# Module 13: Echo Labstack

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

- Understand the core architecture and components of the Echo framework.
- Build robust REST APIs using routing, request handling, and validation.
- Implement middleware for cross-cutting concerns such as logging and authentication.
- Structure an Echo application using best practices for scalability and maintainability.
- Write effective tests for Echo handlers.

## Overview

Go’s standard net/http package provides powerful primitives for building HTTP servers, but building large or complex
APIs directly on top of it can become verbose and repetitive.

Echo is a high-performance, minimalist web framework for Go that simplifies API development while retaining full control
over request handling and error management. Echo provides expressive routing, middleware composition, request binding,
and centralized error handling, making it well suited for building production-ready services.

## Introduction to Echo

Echo is one of the most widely used web frameworks in the Go ecosystem. It is designed to be fast, explicit, and
flexible, enabling developers to build APIs and web applications with a clear and predictable execution model.

Echo emphasizes:

- Explicit handler behavior
- Clear middleware flow
- Centralized error handling
- Minimal abstractions over net/http

For more information, refer to the official documentation:
https://echo.labstack.com/docs

## The Echo Landscape

Before diving into Echo’s specifics, it’s useful to understand what problems it aims to solve.

### Echo and the Standard Library

- Performance: Lightweight routing with minimal overhead
- Developer Experience: Clear handler signatures and structured middleware
- Feature Set: Routing, middleware, binding, validation, and error handling

### Why Echo Matters for Modern Applications

- Simplifies API development
- Encourages explicit control flow
- Scales well with application complexity
- Well suited for RESTful services and microservices

## Syntax Comparison

### Standard `net/http`

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
	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
```

## Core Concepts

### Echo Instance

The Echo instance (*echo.Echo) is the core of the framework. It is responsible for:

- Registering routes
- Managing middleware
- Starting and stopping the HTTP server

```go
e := echo.New()
```

### Context (echo.Context)

The echo.Context represents the context of the current HTTP request. It provides access to:

- Request and response objects
- Path and query parameters
- Request body binding
- Response rendering helpers

Handlers in Echo receive a context and return an error, allowing consistent and centralized error handling.

### Route Groups

Route groups allow related routes to share a common path prefix and middleware chain, commonly used for:

- API versioning
- Feature-based organization
- Secured endpoints

```go
api := e.Group("/api")
```

### Middleware

Middleware are functions that wrap handlers and intercept requests before and/or after the handler executes.
Common use cases include:

- Logging
- Authentication
- Authorization
- Rate limiting
- Error handling

Middleware may choose to continue the request flow or stop it by returning an error.

## Quick Start

### Prerequisite

- Go 1.13 or above

Note: Go 1.12 has limited support and some middlewares will not be available.

### Installation

To install Echo package, you need to install Go and set your Go workspace first. If you don’t have a go.mod file, create
it with `go mod init echo`.

1. Download and install Echo:
    ```shell
    go get github.com/labstack/echo/v4
    ```

2. Create a Simple Server
    ```go
    // file main.go
    package main
    
    import (
        "net/http"
    
        "github.com/labstack/echo/v4"
    )
    
    func main() {
        e := echo.New()
    
        e.GET("/ping", func(c echo.Context) error {
            return c.JSON(http.StatusOK, map[string]string{
                "message": "pong",
            })
        })
    
        e.Logger.Fatal(e.Start(":8080") // listen and serve on 0.0.0.0:8080
    }
    ```
3. Start the server
    ```shell
     go run main.go
    ```   

## Features

### Routing and Handling Requests

Routing maps HTTP methods and paths to handlers. It supports all standard HTTP methods: `GET`, `POST`, `PUT`, `DELETE`,
`PATCH`, `HEAD`, and `OPTIONS`.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/users", saveUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":8080")) // listen and serve on 0.0.0.0:8080
}

func getUser(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"id": c.Param("id"),
	})
}

// Other handler functions (saveUser, updateUser, deleteUser,...)
```

### Working with Route Groups

Organizing related routes into groups improves code structure.

Router Groups are perfect for versioning your API and applying shared logic.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	v1 := e.Group("/api/v1")

	users := v1.Group("/users")
	{
		users.GET("", getUsers)
		users.GET("/:id", getUser)
		users.POST("/:id", createUser)
		users.PUT("/:id", updateUser)
		users.DELETE("/:id", deleteUser)
	}

	e.Logger.Fatal(e.Start(":8080")) // listen and serve on 0.0.0.0:8080
}

// Mock handler functions
func getUser(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "get user detail v1",
	})
}

func getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "get users v1",
	})
}

func createUser(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "create users v1",
	})
}

func updateUser(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "update users v1",
	})
}

func deleteUser(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "delete users v1",
	})
}
```

### Binding payload and parsing data

#### Request

Binding maps incoming request data to Go structs.

```go
package main

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

func createUser(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}
```

#### Path and Query Parameters

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func bindQuery() {
	e := echo.New()
	
	e.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.JSON(http.StatusOK, map[string]string{"id": id})
	})

	e.GET("/users", func(c echo.Context) error {
		page := c.QueryParam("page")
		limit := c.QueryParam("limit")

		return c.JSON(http.StatusOK, map[string]string{
			"page": page,
			"limit":  limit,
		})
	})

	e.POST("/forms", func(c echo.Context) error{
		name := c.FormValue("name")
		return c.String(http.StatusOK, name)
	})
}

```

#### File upload

```go
package main

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
)

func uploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	//Handle upload file
	dst := "uploads/" + file.Filename
	if err := saveFile(file, dst); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"filename": file.Filename,
	})
}

```

### Handling Response

Echo provides helpers for various response types:

```go
c.JSON(http.StatusOK, data)
c.String(http.StatusOK, "plain text")
c.HTML(http.StatusOK, "index.html", data)
c.Redirect(http.StatusMovedPermanently, "/login")
c.File("file.pdf")
```

### Applying Middleware

Middleware can be applied globally or to route groups.

```go
package main

import "github.com/labstack/echo/v4"

func CustomLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Logic before the handler is executed 
		log.Printf("Incoming request: %s %s", c.Request().Method, c.Request().URL.Path)

		// Call the next handler in the chain
		return next(c)
	}
}


e.Use(middleware.Logger())
e.Use(middleware.Recover())
e.Use(CustomLogger)

```

### Chaining Middleware

Middleware execution follows the order in which it is registered. Each middleware explicitly decides whether to continue
the request flow by calling the next handler.

If a middleware returns an error, the chain stops and the error handler is invoked.

A middleware controls the flow by:

- Calling next(c) → continue to the next middleware or handler
- Returning an error → stop the chain and invoke the error handler

### Handling Errors

Echo uses a centralized error handling mechanism.

```go
e.HTTPErrorHandler = func(err error, c echo.Context) {
    code := http.StatusInternalServerError
    if he, ok := err.(*echo.HTTPError); ok {
        code = he.Code
    }
    c.JSON(code, map[string]string{"error": err.Error()})
}
```

Handlers and middleware may return errors to trigger this mechanism.

## Common Echo Patterns and Best Practices

### Project Structure

```
├── main.go          # Entry point
├── config/          # Configuration management
├── controllers/     # HTTP handlers
├── middleware/      # Custom middleware
├── models/          # Data models
├── routes/          # Route definitions
├── services/        # Business logic
├── templates/       # HTML templates
├── utils/           # Helper functions
└── tests/           # Test files

```

### Dependency Injection

```go
package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

// Service interface
type User struct{}

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

func (uc *UserController) GetUser(c echo.Context) error {
	id := c.Param("id")

	user, err := uc.service.GetUser(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, user)
}

func SetupRoutes(r *echo.Echo, uc *UserController) {
	r.GET("/users/:id", uc.GetUser)
	// Other routes
}

func main() {
	e := echo.New()

	// Initialize dependencies
	userService := NewUserService()
	userController := NewUserController(userService)

	// Setup routes
	SetupRoutes(e, userController)

	e.Logger.Fatal(e.Start(":8080"))
}

```

## Common Challenges and Solutions

### Handling CORS

Echo provides middleware support for handling Cross-Origin Resource Sharing (CORS).
CORS middleware allows your API to be accessed safely from browsers running on different origins.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// CORS middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := c.Response()
			req := c.Request()

			res.Header().Set("Access-Control-Allow-Origin", "*")
			res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight request
			if req.Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}

			return next(c)
		}
	})

	// Routes
	// ...

	e.Logger.Fatal(e.Start(":8080"))
}
```

### Rate Limiting

Rate limiting prevents clients from making too many requests in a short period of time, helping protect your application
from abuse.

Below is a simple in-memory rate limiter implemented as Echo middleware.

```go
package main

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

func RateLimiter() echo.MiddlewareFunc {
	// Simple in-memory rate limiter
	limits := make(map[string]int)
	mutex := &sync.Mutex{}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()

			mutex.Lock()
			if limits[ip] >= 100 { // 100 requests per minute
				mutex.Unlock()
				return echo.NewHTTPError(
					http.StatusTooManyRequests,
					"Rate limit exceeded",
				)
			}

			limits[ip]++
			mutex.Unlock()

			// Reset counters periodically (in a real app, use a timer)

			return next(c)
		}
	}
}

```

Applying the Rate Limiter

```go
e.Use(RateLimiter())
```

## Practical Exercises

### Exercise 1: Basic Echo API

Create a simple RESTful API using Echo framework to manage a todo list

### Exercise 2: Echo Middleware and Authentication

Create a Echo application with custom middleware for logging and simple API key authentication

### Exercise 3: File Upload with Echo

Create a Echo application that handles file uploads with progress monitoring