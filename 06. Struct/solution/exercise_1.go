package main

import (
	"fmt"
	"time"
)

// Book represents a book in the library
type Book struct {
	ID            string
	Title         string
	Author        string
	PublishedYear int
	Available     bool
}

// Member represents a library member
type Member struct {
	ID       string
	Name     string
	Email    string
	JoinedOn time.Time
	BooksOut int
	MaxBooks int
}

// BorrowRecord tracks a book being borrowed
type BorrowRecord struct {
	BookID     string
	MemberID   string
	BorrowedOn time.Time
	DueDate    time.Time
	ReturnedOn *time.Time // Pointer because it might be nil (not returned yet)
}

// Library manages the book collection and members
type Library struct {
	Name    string
	Books   map[string]*Book
	Members map[string]*Member
	Borrows []BorrowRecord
}

// NewLibrary creates a new library instance
func NewLibrary(name string) *Library {
	return &Library{
		Name:    name,
		Books:   make(map[string]*Book),
		Members: make(map[string]*Member),
		Borrows: []BorrowRecord{},
	}
}

// AddBook adds a book to the library
func (l *Library) AddBook(book Book) {
	l.Books[book.ID] = &book
}

// AddMember adds a member to the library
func (l *Library) AddMember(member Member) {
	l.Members[member.ID] = &member
}

// BorrowBook allows a member to borrow a book
func (l *Library) BorrowBook(bookID, memberID string) error {
	// Find the book
	book, found := l.Books[bookID]
	if !found {
		return fmt.Errorf("book not found")
	}

	// Check if the book is available
	if !book.Available {
		return fmt.Errorf("book is not available")
	}

	// Find the member
	member, found := l.Members[memberID]
	if !found {
		return fmt.Errorf("member not found")
	}

	// Check if the member can borrow more books
	if member.BooksOut >= member.MaxBooks {
		return fmt.Errorf("member has reached maximum number of books")
	}

	// Create a borrow record
	now := time.Now()
	borrowRecord := BorrowRecord{
		BookID:     bookID,
		MemberID:   memberID,
		BorrowedOn: now,
		DueDate:    now.AddDate(0, 0, 14), // Due in 14 days
	}

	// Update book and member
	book.Available = false
	member.BooksOut++

	// Add the record
	l.Borrows = append(l.Borrows, borrowRecord)

	return nil
}

// ReturnBook processes a book return
func (l *Library) ReturnBook(bookID, memberID string) error {
	// Find the book
	book, found := l.Books[bookID]
	if !found {
		return fmt.Errorf("book not found")
	}

	// Find the member
	member, found := l.Members[memberID]
	if !found {
		return fmt.Errorf("member not found")
	}

	// Find the borrow record
	recordIndex := -1
	for i, record := range l.Borrows {
		if record.BookID == bookID && record.MemberID == memberID && record.ReturnedOn == nil {
			recordIndex = i
			break
		}
	}

	if recordIndex == -1 {
		return fmt.Errorf("no active borrow record found")
	}

	// Update the record
	now := time.Now()
	l.Borrows[recordIndex].ReturnedOn = &now

	// Update book and member
	book.Available = true
	member.BooksOut--

	return nil
}

func main() {
	// Create a new library
	library := NewLibrary("Community Library")

	// Add books
	library.AddBook(Book{
		ID:            "B001",
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan & Brian W. Kernighan",
		PublishedYear: 2015,
		Available:     true,
	})

	library.AddBook(Book{
		ID:            "B002",
		Title:         "Go in Action",
		Author:        "William Kennedy",
		PublishedYear: 2016,
		Available:     true,
	})

	// Add members
	library.AddMember(Member{
		ID:       "M001",
		Name:     "John Doe",
		Email:    "john@example.com",
		JoinedOn: time.Now(),
		BooksOut: 0,
		MaxBooks: 3,
	})

	library.AddMember(Member{
		ID:       "M002",
		Name:     "Jane Smith",
		Email:    "jane@example.com",
		JoinedOn: time.Now(),
		BooksOut: 0,
		MaxBooks: 5,
	})

	// Borrow a book
	err := library.BorrowBook("B001", "M001")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Println("Book B001 borrowed by member M001")
	}

	// Display library status
	fmt.Println("\nLibrary Status:")
	fmt.Printf("Name: %s\n", library.Name)
	fmt.Printf("Books: %d\n", len(library.Books))
	fmt.Printf("Members: %d\n", len(library.Members))
	fmt.Printf("Active Borrows: %d\n", len(library.Borrows))

	// Return the book
	err = library.ReturnBook("B001", "M001")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Println("\nBook B001 returned by member M001")
	}
}
