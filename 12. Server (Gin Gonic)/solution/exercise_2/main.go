package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
func CustomLogger() gin.HandlerFunc {
	// Create log file
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s | %d | %s | %s | %s\n",
			param.TimeStamp.Format(time.RFC822),
			param.ClientIP,
			param.StatusCode,
			param.Method,
			param.Path,
			param.Latency,
		)
	})
}

// APIKeyAuth implements authentication using API keys
func APIKeyAuth(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get API key from header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// Check if it's in query string
			apiKey = c.Query("api_key")
		}

		// Validate API key
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key is required",
			})
			return
		}

		username, valid := config.APIKeys[apiKey]
		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			return
		}

		// Store user information in the context
		c.Set("user", username)
		c.Next()
	}
}

// GetUserFromContext retrieves the user from the Gin context
func GetUserFromContext(c *gin.Context) string {
	user, exists := c.Get("user")
	if !exists {
		return "Unknown"
	}
	return user.(string)
}

func main() {
	config := NewConfig()

	// Create a Gin router with default middleware
	r := gin.New()

	// Add custom middlewares
	r.Use(CustomLogger())
	r.Use(gin.Recovery())

	// Public endpoints
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Secure API",
			"status":  "online",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Secured API group
	api := r.Group("/api")
	api.Use(APIKeyAuth(config))
	{
		api.GET("/protected", func(c *gin.Context) {
			username := GetUserFromContext(c)
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Hello, %s! This is protected data.", username),
				"time":    time.Now().Format(time.RFC3339),
			})
		})

		api.GET("/profile", func(c *gin.Context) {
			username := GetUserFromContext(c)
			role := "user"
			if username == "Administrator" {
				role = "admin"
			}

			c.JSON(http.StatusOK, gin.H{
				"username": username,
				"role":     role,
				"access":   "granted",
			})
		})
	}

	// Start the server
	log.Println("Starting secure API server on :8080...")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
