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
services while keeping code maintainable and efficient. This course will guide you step by step, from the basics to
advanced practices,
so you can master Echo and apply it to real-world projects.

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

The Echo Engine is the core of the framework. It’s the main instance you create with echo.New(). The engine is
responsible
for registering routes, attaching middleware, and starting the HTTP server. You can customize it by adding middleware,
defining
handlers, and configuring server settings.

### Context (echo.Context)

The `Context` is the most important component in Echo. It encapsulates the HTTP request and response, giving you
convenient
methods to read parameters, headers, query strings, and request bodies. It also provides functions to send responses in
various
formats (JSON, HTML, text, file, etc.), making request–response handling simple and efficient.

### Router Groups

Echo provides Group functionality to organize routes under a shared path prefix and middleware. For example, you can
group routes
by version (e.g., `/api/v1/users`, `/api/v1/products`) and apply authentication middleware only once at the group level.
This keeps
code modular and easier to maintain.

### Middleware

Middleware in Echo are functions executed before or after handlers. They can be applied globally, per route, or per
group. Common use
cases include logging, error handling, authentication, authorization, rate limiting, and CORS. Middleware can control
request flow by
calling `next(c)` to continue or aborting with an error response.

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
    
         e.Start(":8080") // listen and serve on 0.0.0.0:8080
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
    
        e.Start(":8080") // listen and serve on 0.0.0.0:8080
    }
```

Echo provides an intuitive API for defining routes with different HTTP methods,
supports all standard HTTP methods: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, and `OPTIONS`.

```go
    package main
    
    import (
        "github.com/labstack/echo/v4"
        "net/http"
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