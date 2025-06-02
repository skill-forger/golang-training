## Practical Exercises

### Exercise 1: Basic Gin API

Create a simple RESTful API using Gin framework to manage a todo list:

```go
// todo_api.go
package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Todo represents a todo item
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" binding:"required"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TodoStore manages the todo items
type TodoStore struct {
	todos  []Todo
	nextID int
}

// NewTodoStore creates a new store with initial data
func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos: []Todo{
			{
				ID:        1,
				Title:     "Learn Gin Framework",
				Completed: false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        2,
				Title:     "Build a RESTful API",
				Completed: false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		nextID: 3,
	}
}

func main() {
	store := NewTodoStore()

	// Create a default gin router
	r := gin.Default()

	// Define API routes
	v1 := r.Group("/api/v1")
	{
		// GET /api/v1/todos - Get all todos
		v1.GET("/todos", func(c *gin.Context) {
			c.JSON(http.StatusOK, store.todos)
		})

		// GET /api/v1/todos/:id - Get a specific todo
		v1.GET("/todos/:id", func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
				return
			}

			// Find the todo
			for _, todo := range store.todos {
				if todo.ID == id {
					c.JSON(http.StatusOK, todo)
					return
				}
			}

			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		})

		// POST /api/v1/todos - Create a new todo
		v1.POST("/todos", func(c *gin.Context) {
			var newTodo Todo

			// Bind JSON body to the newTodo struct
			if err := c.ShouldBindJSON(&newTodo); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Set todo properties
			newTodo.ID = store.nextID
			store.nextID++
			newTodo.Completed = false
			newTodo.CreatedAt = time.Now()
			newTodo.UpdatedAt = time.Now()

			// Add to store
			store.todos = append(store.todos, newTodo)

			c.JSON(http.StatusCreated, newTodo)
		})

		// PUT /api/v1/todos/:id - Update a todo
		v1.PUT("/todos/:id", func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
				return
			}

			var updatedTodo Todo
			if err := c.ShouldBindJSON(&updatedTodo); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			for i, todo := range store.todos {
				if todo.ID == id {
					// Preserve ID and creation time
					updatedTodo.ID = id
					updatedTodo.CreatedAt = todo.CreatedAt
					updatedTodo.UpdatedAt = time.Now()

					store.todos[i] = updatedTodo
					c.JSON(http.StatusOK, updatedTodo)
					return
				}
			}

			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		})

		// DELETE /api/v1/todos/:id - Delete a todo
		v1.DELETE("/todos/:id", func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
				return
			}

			for i, todo := range store.todos {
				if todo.ID == id {
					// Remove the todo
					store.todos = append(store.todos[:i], store.todos[i+1:]...)
					c.Status(http.StatusNoContent)
					return
				}
			}

			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		})
	}

	// Start the server
	r.Run(":8080")
}
```

### Exercise 2: Gin Middleware and Authentication

Create a Gin application with custom middleware for logging and simple API key authentication:

```go
// secure_api.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Config holds the application configuration
type Config struct {
	APIKeys map[string]string // Map of API key to user name
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
			c.JSON(http.StatusOK, gin.H{
				"username": username,
				"role":     username == "Administrator" ? "admin" : "user",
				"access":   "granted",
			})
		})
	}

	// Start the server
	log.Println("Starting secure API server on :8080...")
	r.Run(":8080")
}
```

### Exercise 3: File Upload with Gin

Create a Gin application that handles file uploads with progress monitoring:

```go
// file_upload.go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Upload statistics
type UploadStats struct {
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// In-memory store for upload stats
var uploads []UploadStats

// Maximum file size (10 MB)
const maxFileSize = 10 * 1024 * 1024

func main() {
	// Create uploads directory if it doesn't exist
	os.MkdirAll("./uploads", 0755)

	// Create a Gin router with default middleware
	r := gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Serve static files from the uploads directory
	r.Static("/files", "./uploads")

	// Serve the HTML upload form
	r.GET("/", func(c *gin.Context) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Go File Upload</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
				h1 { color: #333; }
				.upload-form { margin: 20px 0; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
				.progress { width: 100%; background-color: #f3f3f3; border-radius: 5px; }
				.progress-bar { height: 20px; background-color: #4CAF50; border-radius: 5px; width: 0%; transition: width 0.3s; }
				.file-item { margin: 10px 0; padding: 10px; border: 1px solid #eee; border-radius: 5px; }
			</style>
		</head>
		<body>
			<h1>File Upload with Go & Gin</h1>
			
			<div class="upload-form">
				<h2>Upload File</h2>
				<form id="uploadForm" enctype="multipart/form-data">
					<input type="file" name="file" required>
					<button type="submit">Upload</button>
				</form>
				<div class="progress" style="display:none;">
					<div class="progress-bar" id="progressBar"></div>
					<div id="progressText">0%</div>
				</div>
			</div>
			
			<h2>Uploaded Files</h2>
			<div id="fileList">
				Loading files...
			</div>
			
			<script>
				// Load file list on page load
				document.addEventListener('DOMContentLoaded', loadFiles);
				
				// Handle form submission
				document.getElementById('uploadForm').addEventListener('submit', function(e) {
					e.preventDefault();
					
					const fileInput = document.querySelector('input[name="file"]');
					const file = fileInput.files[0];
					
					if (!file) {
						alert('Please select a file to upload');
						return;
					}
					
					const formData = new FormData();
					formData.append('file', file);
					
					const xhr = new XMLHttpRequest();
					
					// Show progress bar
					const progressBar = document.getElementById('progressBar');
					const progressText = document.getElementById('progressText');
					document.querySelector('.progress').style.display = 'block';
					
					// Track upload progress
					xhr.upload.addEventListener('progress', function(e) {
						if (e.lengthComputable) {
							const percentComplete = Math.round((e.loaded / e.total) * 100);
							progressBar.style.width = percentComplete + '%';
							progressText.textContent = percentComplete + '%';
						}
					});
					
					xhr.addEventListener('load', function() {
						if (xhr.status === 200) {
							alert('File uploaded successfully!');
							fileInput.value = '';
							loadFiles();
						} else {
							alert('Upload failed: ' + xhr.responseText);
						}
						
						// Hide progress bar after upload
						setTimeout(() => {
							document.querySelector('.progress').style.display = 'none';
							progressBar.style.width = '0%';
							progressText.textContent = '0%';
						}, 1000);
					});
					
					xhr.open('POST', '/upload', true);
					xhr.send(formData);
				});
				
				// Load list of uploaded files
				function loadFiles() {
					fetch('/files-list')
						.then(response => response.json())
						.then(data => {
							const fileList = document.getElementById('fileList');
							
							if (data.length === 0) {
								fileList.innerHTML = '<p>No files uploaded yet.</p>';
								return;
							}
							
							let html = '';
							data.forEach(file => {
								const sizeInMB = (file.size / (1024 * 1024)).toFixed(2);
								html += `
									<div class="file-item">
										<strong>${file.filename}</strong> (${sizeInMB} MB)
										<div>Type: ${file.mime_type}</div>
										<div>Uploaded: ${new Date(file.uploaded_at).toLocaleString()}</div>
										<a href="/files/${file.filename}" target="_blank">Download</a>
									</div>
								`;
							});
							
							fileList.innerHTML = html;
						})
						.catch(error => {
							console.error('Error loading files:', error);
							document.getElementById('fileList').innerHTML = '<p>Error loading files.</p>';
						});
				}
			</script>
		</body>
		</html>
		`;
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	})

	// Handle file upload
	r.POST("/upload", func(c *gin.Context) {
		// Get the file from the request
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		// Check file size
		if header.Size > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("File too large (max %d MB)", maxFileSize/(1024*1024)),
			})
			return
		}

		// Create a safe filename
		filename := header.Filename
		// Remove any path from the filename
		filename = filepath.Base(filename)
		// Ensure filename is unique
		ext := filepath.Ext(filename)
		basename := strings.TrimSuffix(filename, ext)
		filename = fmt.Sprintf("%s_%d%s", basename, time.Now().Unix(), ext)

		// Create destination file
		dst, err := os.Create(filepath.Join("uploads", filename))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer dst.Close()

		// Copy the file
		_, err = io.Copy(dst, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Detect mime type
		mimeType := header.Header.Get("Content-Type")
		
		// Store upload stats
		stats := UploadStats{
			Filename:   filename,
			Size:       header.Size,
			MimeType:   mimeType,
			UploadedAt: time.Now(),
		}
		uploads = append(uploads, stats)

		c.JSON(http.StatusOK, stats)
	})

	// Get list of uploaded files
	r.GET("/files-list", func(c *gin.Context) {
		c.JSON(http.StatusOK, uploads)
	})

	// Start the server
	log.Println("Starting file upload server on :8080...")
	r.Run(":8080")
}
