package main

import (
	"log"
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
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
