package main

import (
	"fmt"
	"sort"
	"strings"
)

// Contact holds information about a person
type Contact struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// ContactBook manages a collection of contacts
type ContactBook struct {
	contacts map[string]Contact
}

// NewContactBook creates a new contact book
func NewContactBook() *ContactBook {
	return &ContactBook{
		contacts: make(map[string]Contact),
	}
}

// AddContact adds a new contact to the book
func (cb *ContactBook) AddContact(contact Contact) {
	key := getContactKey(contact)
	cb.contacts[key] = contact
}

// getContactKey creates a unique key for a contact
func getContactKey(c Contact) string {
	return strings.ToLower(c.FirstName + ":" + c.LastName)
}

// FindContact searches for contacts by name
func (cb *ContactBook) FindContact(name string) []Contact {
	var results []Contact
	name = strings.ToLower(name)

	for _, contact := range cb.contacts {
		firstName := strings.ToLower(contact.FirstName)
		lastName := strings.ToLower(contact.LastName)

		if strings.Contains(firstName, name) || strings.Contains(lastName, name) {
			results = append(results, contact)
		}
	}

	// Sort results by last name, then first name
	sort.Slice(results, func(i, j int) bool {
		if results[i].LastName == results[j].LastName {
			return results[i].FirstName < results[j].FirstName
		}
		return results[i].LastName < results[j].LastName
	})

	return results
}

// DeleteContact removes a contact by their full name
func (cb *ContactBook) DeleteContact(firstName, lastName string) bool {
	key := strings.ToLower(firstName + ":" + lastName)
	_, exists := cb.contacts[key]
	if exists {
		delete(cb.contacts, key)
		return true
	}
	return false
}

// ListAllContacts returns all contacts in alphabetical order
func (cb *ContactBook) ListAllContacts() []Contact {
	var allContacts []Contact
	for _, contact := range cb.contacts {
		allContacts = append(allContacts, contact)
	}

	// Sort by last name, then first name
	sort.Slice(allContacts, func(i, j int) bool {
		if allContacts[i].LastName == allContacts[j].LastName {
			return allContacts[i].FirstName < allContacts[j].FirstName
		}
		return allContacts[i].LastName < allContacts[j].LastName
	})

	return allContacts
}

func main() {
	// Create a new contact book
	book := NewContactBook()

	// Add some sample contacts
	book.AddContact(Contact{"John", "Doe", "john.doe@example.com", "555-1234"})
	book.AddContact(Contact{"Jane", "Smith", "jane.smith@example.com", "555-5678"})
	book.AddContact(Contact{"Alice", "Johnson", "alice.j@example.com", "555-9012"})
	book.AddContact(Contact{"Bob", "Brown", "bob.brown@example.com", "555-3456"})
	book.AddContact(Contact{"John", "Smith", "john.smith@example.com", "555-7890"})

	// List all contacts
	fmt.Println("All Contacts:")
	fmt.Println("-------------")
	for i, contact := range book.ListAllContacts() {
		fmt.Printf("%d. %s %s\n   Email: %s\n   Phone: %s\n\n",
			i+1, contact.FirstName, contact.LastName,
			contact.Email, contact.Phone)
	}

	// Search for contacts
	searchTerm := "john"
	fmt.Printf("\nSearch results for '%s':\n", searchTerm)
	fmt.Println("---------------------------")
	results := book.FindContact(searchTerm)
	if len(results) == 0 {
		fmt.Println("No contacts found.")
	} else {
		for i, contact := range results {
			fmt.Printf("%d. %s %s\n   Email: %s\n   Phone: %s\n\n",
				i+1, contact.FirstName, contact.LastName,
				contact.Email, contact.Phone)
		}
	}

	// Delete a contact
	deleted := book.DeleteContact("John", "Doe")
	fmt.Printf("\nDeleted John Doe: %t\n", deleted)

	// List contacts after deletion
	fmt.Println("\nRemaining Contacts:")
	fmt.Println("------------------")
	for i, contact := range book.ListAllContacts() {
		fmt.Printf("%d. %s %s\n", i+1, contact.FirstName, contact.LastName)
	}
}
