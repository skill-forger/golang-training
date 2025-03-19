## Practical Exercises

### Exercise 1: Student Grade Tracker

Create a program that tracks and analyzes student grades using maps and slices. This exercise will help you understand how to work with collection types to store and process related data.

Your implementation should:
1. Create a map that associates student names (strings) with their grades (slices of integers)
2. Calculate each student's average grade and store it in a new map
3. Identify the student with the highest average grade
4. Create a ranked list of students based on their average grades
5. Display formatted output showing:
   - Individual student averages
   - The top performing student
   - A ranked list of all students

```go
// student_grades.go
package main

import (
    "fmt"
    "sort"
)

func main() {
    // Initialize student grades
    grades := map[string][]int{
        "Alice":   {92, 88, 95, 89},
        "Bob":     {75, 82, 79},
        "Charlie": {90, 93, 88, 97, 91},
        "Diana":   {65, 72, 80, 75},
    }
    
    // Calculate average grades
    averages := make(map[string]float64)
    for student, studentGrades := range grades {
        total := 0
        for _, grade := range studentGrades {
            total += grade
        }
        averages[student] = float64(total) / float64(len(studentGrades))
    }
    
    // Print averages
    fmt.Println("Student Average Grades:")
    for student, avg := range averages {
        fmt.Printf("%s: %.2f\n", student, avg)
    }
    
    // Find the student with the highest average
    var topStudent string
    var topAverage float64
    
    for student, avg := range averages {
        if avg > topAverage {
            topAverage = avg
            topStudent = student
        }
    }
    
    fmt.Printf("\nTop student: %s with average %.2f\n", topStudent, topAverage)
    
    // List all students ordered by grade (highest first)
    type StudentAvg struct {
        Name    string
        Average float64
    }
    
    var studentList []StudentAvg
    for name, avg := range averages {
        studentList = append(studentList, StudentAvg{name, avg})
    }
    
    // Sort by average (descending)
    sort.Slice(studentList, func(i, j int) bool {
        return studentList[i].Average > studentList[j].Average
    })
    
    fmt.Println("\nRanked Students:")
    for i, s := range studentList {
        fmt.Printf("%d. %s: %.2f\n", i+1, s.Name, s.Average)
    }
}
```

### Exercise 2: Word Frequency Counter

Develop a text analysis tool that counts word frequencies in a text document. This exercise demonstrates how to process strings, use regular expressions, and manipulate maps for data analysis.

Your implementation should:
1. Process a sample text (or optionally read from a file)
2. Convert the text to lowercase and extract individual words
3. Count the frequency of each word in the text using a map
4. Filter out common stop words that don't add meaning
5. Sort the words by frequency in descending order
6. Display the top N most frequent words in a formatted table
7. Optionally write the complete results to a file

```go
// word_counter.go
package main

import (
    "fmt"
    "os"
    "regexp"
    "sort"
    "strings"
)

func main() {
    // Sample text or you could read from a file
    text := `Go is an open source programming language that makes it easy to build 
    simple, reliable, and efficient software. Go was designed at Google in 2007 
    by Robert Griesemer, Rob Pike, and Ken Thompson. Go is syntactically similar 
    to C, but with memory safety, garbage collection, structural typing, and 
    CSP-style concurrency. The language is often referred to as Golang because of 
    its former domain name, golang.org, but the proper name is Go.`
    
    // Convert to lowercase and split into words
    text = strings.ToLower(text)
    
    // Use regex to extract words
    re := regexp.MustCompile(`[a-z]+`)
    words := re.FindAllString(text, -1)
    
    // Count word frequencies
    frequencies := make(map[string]int)
    for _, word := range words {
        frequencies[word]++
    }
    
    // Filter out common words (simplified stop words list)
    stopWords := map[string]bool{
        "the": true, "and": true, "is": true, "to": true, "of": true, 
        "a": true, "in": true, "but": true, "with": true, "by": true,
        "was": true, "its": true,
    }
    
    // Create a list of word-frequency pairs
    type WordFreq struct {
        Word  string
        Count int
    }
    
    var wordFreqs []WordFreq
    for word, count := range frequencies {
        if !stopWords[word] && len(word) > 1 {
            wordFreqs = append(wordFreqs, WordFreq{word, count})
        }
    }
    
    // Sort by frequency (descending)
    sort.Slice(wordFreqs, func(i, j int) bool {
        return wordFreqs[i].Count > wordFreqs[j].Count
    })
    
    // Print top N words
    topN := 10
    if len(wordFreqs) < topN {
        topN = len(wordFreqs)
    }
    
    fmt.Printf("Top %d words:\n", topN)
    fmt.Printf("%-15s %s\n", "WORD", "FREQUENCY")
    fmt.Println("------------------------")
    
    for i := 0; i < topN; i++ {
        wf := wordFreqs[i]
        fmt.Printf("%-15s %d\n", wf.Word, wf.Count)
    }
    
    // Print to file (optional)
    file, err := os.Create("word_frequencies.txt")
    if err == nil {
        defer file.Close()
        
        fmt.Fprintf(file, "%-15s %s\n", "WORD", "FREQUENCY")
        fmt.Fprintln(file, "------------------------")
        
        for _, wf := range wordFreqs {
            fmt.Fprintf(file, "%-15s %d\n", wf.Word, wf.Count)
        }
    }
}
```

### Exercise 3: Contact Book Application

Build a contact management application that allows users to store and retrieve contact information. This exercise combines structs with collection types to create a more complex data management system.

Your implementation should include:
1. A `Contact` struct that stores personal information (first name, last name, email, phone)
2. A `ContactBook` struct that manages a collection of contacts using a map
3. Methods for:
   - Adding a new contact to the book
   - Finding contacts by searching for a name or partial name
   - Deleting contacts
   - Listing all contacts in alphabetical order
   - Grouping contacts by their first letter
4. A demonstration that shows all the functionality of the contact book
5. Proper handling of case sensitivity in searches
6. Sorting capabilities for displaying contacts in a structured way

```go
// contact_book.go
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
```
