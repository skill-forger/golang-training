package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// Todo represents a todo item
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
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
				Title:     "Learn Echo Framework",
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

	// Create Echo instance
	e := echo.New()

	// API version group
	v1 := e.Group("/api/v1")

	// GET /api/v1/todos - Get all todos
	v1.GET("/todos", func(c echo.Context) error {
		return c.JSON(http.StatusOK, store.todos)
	})

	// GET /api/v1/todos/:id - Get a specific todo
	v1.GET("/todos/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid todo ID")
		}

		for _, todo := range store.todos {
			if todo.ID == id {
				return c.JSON(http.StatusOK, todo)
			}
		}

		return echo.NewHTTPError(http.StatusNotFound, "Todo not found")
	})

	// POST /api/v1/todos - Create a new todo
	v1.POST("/todos", func(c echo.Context) error {
		var newTodo Todo

		if err := c.Bind(&newTodo); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		newTodo.ID = store.nextID
		store.nextID++
		newTodo.Completed = false
		newTodo.CreatedAt = time.Now()
		newTodo.UpdatedAt = time.Now()

		store.todos = append(store.todos, newTodo)

		return c.JSON(http.StatusCreated, newTodo)
	})

	// PUT /api/v1/todos/:id - Update a todo
	v1.PUT("/todos/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid todo ID")
		}

		var updatedTodo Todo
		if err := c.Bind(&updatedTodo); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		for i, todo := range store.todos {
			if todo.ID == id {
				updatedTodo.ID = id
				updatedTodo.CreatedAt = todo.CreatedAt
				updatedTodo.UpdatedAt = time.Now()

				store.todos[i] = updatedTodo
				return c.JSON(http.StatusOK, updatedTodo)
			}
		}

		return echo.NewHTTPError(http.StatusNotFound, "Todo not found")
	})

	// DELETE /api/v1/todos/:id - Delete a todo
	v1.DELETE("/todos/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid todo ID")
		}

		for i, todo := range store.todos {
			if todo.ID == id {
				store.todos = append(store.todos[:i], store.todos[i+1:]...)
				return c.NoContent(http.StatusNoContent)
			}
		}

		return echo.NewHTTPError(http.StatusNotFound, "Todo not found")
	})

	// Start server
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
