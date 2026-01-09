package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Config holds the application configuration
type Config struct {
	APIKeys map[string]string // Map of API key to username
}

// NewConfig creates a default configuration
func NewConfig() *Config {
	return &Config{
		APIKeys: map[string]string{
			"development-key": "Developer",
			"test-key":        "Tester",
			"admin-key":       "Administrator",
		},
	}
}

// CustomLogger implements a custom logging middleware
func CustomLogger() echo.MiddlewareFunc {
	// Create log file
	f, _ := os.Create("echo.log")
	writer := io.MultiWriter(f, os.Stdout)
	logger := log.New(writer, "", 0)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			req := c.Request()
			res := c.Response()

			logger.Printf(
				"[%s] %s | %d | %s | %s | %s",
				time.Now().Format(time.RFC822),
				c.RealIP(),
				res.Status,
				req.Method,
				req.URL.Path,
				time.Since(start),
			)

			return err
		}
	}
}

// APIKeyAuth implements authentication using API keys
func APIKeyAuth(config *Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get API key from header
			apiKey := c.Request().Header.Get("X-API-Key")
			if apiKey == "" {
				// Fallback to query parameter
				apiKey = c.QueryParam("api_key")
			}

			if apiKey == "" {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					"API key is required",
				)
			}

			username, valid := config.APIKeys[apiKey]
			if !valid {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					"Invalid API key",
				)
			}

			// Store user information in the context
			c.Set("user", username)

			return next(c)
		}
	}
}

// GetUserFromContext retrieves the user from the Echo context
func GetUserFromContext(c echo.Context) string {
	user := c.Get("user")
	if user == nil {
		return "Unknown"
	}
	return user.(string)
}

func main() {
	config := NewConfig()

	// Create Echo instance
	e := echo.New()

	// Register custom middlewares
	e.Use(CustomLogger())
	e.Use(middleware.Recover())

	// Public endpoint
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to the Secure API",
			"status":  "online",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Secured API group
	api := e.Group("/api")
	api.Use(APIKeyAuth(config))

	api.GET("/protected", func(c echo.Context) error {
		username := GetUserFromContext(c)
		return c.JSON(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Hello, %s! This is protected data.", username),
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	api.GET("/profile", func(c echo.Context) error {
		username := GetUserFromContext(c)
		role := "user"
		if username == "Administrator" {
			role = "admin"
		}

		return c.JSON(http.StatusOK, map[string]string{
			"username": username,
			"role":     role,
			"access":   "granted",
		})
	})

	// Start server
	log.Println("Starting secure API server on :8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
