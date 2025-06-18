package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Book represents a book entity
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

// BookStore manages the collection of books
type BookStore struct {
	books  []Book
	nextID int
}

// NewBookStore creates a new book store with some initial data
func NewBookStore() *BookStore {
	return &BookStore{
		books: []Book{
			{ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan & Brian Kernighan", Year: 2015},
			{ID: 2, Title: "Go in Action", Author: "William Kennedy", Year: 2016},
		},
		nextID: 3,
	}
}

func main() {
	store := NewBookStore()

	// Define handlers
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// GET /books - Return all books
			handleGetBooks(w, store)
		case http.MethodPost:
			// POST /books - Create a new book
			handleCreateBook(w, r, store)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		// Extract book ID from URL
		idStr := strings.TrimPrefix(r.URL.Path, "/books/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid book ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// GET /books/{id} - Get a specific book
			handleGetBook(w, id, store)
		case http.MethodPut:
			// PUT /books/{id} - Update a specific book
			handleUpdateBook(w, r, id, store)
		case http.MethodDelete:
			// DELETE /books/{id} - Delete a specific book
			handleDeleteBook(w, id, store)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start server
	fmt.Println("Starting book server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler functions

func handleGetBooks(w http.ResponseWriter, store *BookStore) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.books)
}

func handleCreateBook(w http.ResponseWriter, r *http.Request, store *BookStore) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Assign a new ID
	book.ID = store.nextID
	store.nextID++

	// Add to collection
	store.books = append(store.books, book)

	// Return the created book
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func handleGetBook(w http.ResponseWriter, id int, store *BookStore) {
	for _, book := range store.books {
		if book.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

func handleUpdateBook(w http.ResponseWriter, r *http.Request, id int, store *BookStore) {
	var updatedBook Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, book := range store.books {
		if book.ID == id {
			// Preserve the book ID
			updatedBook.ID = id
			store.books[i] = updatedBook

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

func handleDeleteBook(w http.ResponseWriter, id int, store *BookStore) {
	for i, book := range store.books {
		if book.ID == id {
			// Remove the book
			store.books = append(store.books[:i], store.books[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}
