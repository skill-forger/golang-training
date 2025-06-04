package main

import (
	"fmt"
	"time"
)

// Address represents a physical address
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

// Employee defines the base employee structure
type Employee struct {
	ID         string
	FirstName  string
	LastName   string
	Email      string
	HireDate   time.Time
	Address    Address
	Position   string
	Salary     float64
	ManagerID  string
	Department string
	IsActive   bool
}

// FullName returns the employee's full name
func (e Employee) FullName() string {
	return e.FirstName + " " + e.LastName
}

// YearsOfService calculates the years an employee has worked
func (e Employee) YearsOfService() float64 {
	now := time.Now()
	duration := now.Sub(e.HireDate)
	return duration.Hours() / 24 / 365.25
}

// Company contains all employees and departments
type Company struct {
	Name        string
	Employees   map[string]*Employee
	Departments map[string][]string // Department name -> slice of employee IDs
}

// NewCompany creates a new company
func NewCompany(name string) *Company {
	return &Company{
		Name:        name,
		Employees:   make(map[string]*Employee),
		Departments: make(map[string][]string),
	}
}

// AddEmployee adds a new employee to the company
func (c *Company) AddEmployee(e Employee) error {
	// Check if employee already exists
	if _, exists := c.Employees[e.ID]; exists {
		return fmt.Errorf("employee with ID %s already exists", e.ID)
	}

	// Add employee to the map
	c.Employees[e.ID] = &e

	// Add employee to the department
	if e.Department != "" {
		c.Departments[e.Department] = append(c.Departments[e.Department], e.ID)
	}

	return nil
}

// UpdateSalary updates an employee's salary
func (c *Company) UpdateSalary(employeeID string, newSalary float64) error {
	employee, exists := c.Employees[employeeID]
	if !exists {
		return fmt.Errorf("employee with ID %s not found", employeeID)
	}

	employee.Salary = newSalary
	return nil
}

// TransferEmployee moves an employee to a new department
func (c *Company) TransferEmployee(employeeID, newDepartment string) error {
	employee, exists := c.Employees[employeeID]
	if !exists {
		return fmt.Errorf("employee with ID %s not found", employeeID)
	}

	// Remove from old department
	oldDepartment := employee.Department
	if oldDepartment != "" {
		for i, id := range c.Departments[oldDepartment] {
			if id == employeeID {
				// Remove by replacing with last element and truncating
				c.Departments[oldDepartment][i] = c.Departments[oldDepartment][len(c.Departments[oldDepartment])-1]
				c.Departments[oldDepartment] = c.Departments[oldDepartment][:len(c.Departments[oldDepartment])-1]
				break
			}
		}
	}

	// Add to new department
	employee.Department = newDepartment
	c.Departments[newDepartment] = append(c.Departments[newDepartment], employeeID)

	return nil
}

// TerminateEmployee marks an employee as inactive
func (c *Company) TerminateEmployee(employeeID string) error {
	employee, exists := c.Employees[employeeID]
	if !exists {
		return fmt.Errorf("employee with ID %s not found", employeeID)
	}

	employee.IsActive = false
	return nil
}

// GetDepartmentStats returns statistics about a department
func (c *Company) GetDepartmentStats(department string) (count int, avgSalary float64, avgService float64) {
	employeeIDs, exists := c.Departments[department]
	if !exists || len(employeeIDs) == 0 {
		return 0, 0, 0
	}

	activeCount := 0
	totalSalary := 0.0
	totalService := 0.0

	for _, id := range employeeIDs {
		employee := c.Employees[id]
		if employee.IsActive {
			activeCount++
			totalSalary += employee.Salary
			totalService += employee.YearsOfService()
		}
	}

	if activeCount == 0 {
		return 0, 0, 0
	}

	return activeCount, totalSalary / float64(activeCount), totalService / float64(activeCount)
}

func main() {
	// Create a new company
	company := NewCompany("Tech Innovations Inc.")

	// Add employees
	employees := []Employee{
		{
			ID:        "E001",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			HireDate:  time.Date(2018, 5, 15, 0, 0, 0, 0, time.UTC),
			Address: Address{
				Street:     "123 Main St",
				City:       "Boston",
				State:      "MA",
				PostalCode: "02108",
				Country:    "USA",
			},
			Position:   "Software Engineer",
			Salary:     95000,
			Department: "Engineering",
			IsActive:   true,
		},
		{
			ID:        "E002",
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			HireDate:  time.Date(2019, 3, 10, 0, 0, 0, 0, time.UTC),
			Address: Address{
				Street:     "456 Park Ave",
				City:       "New York",
				State:      "NY",
				PostalCode: "10022",
				Country:    "USA",
			},
			Position:   "Marketing Specialist",
			Salary:     85000,
			Department: "Marketing",
			IsActive:   true,
		},
		{
			ID:        "E003",
			FirstName: "Robert",
			LastName:  "Johnson",
			Email:     "robert.johnson@example.com",
			HireDate:  time.Date(2017, 1, 20, 0, 0, 0, 0, time.UTC),
			Address: Address{
				Street:     "789 Oak St",
				City:       "San Francisco",
				State:      "CA",
				PostalCode: "94107",
				Country:    "USA",
			},
			Position:   "Senior Software Engineer",
			Salary:     120000,
			Department: "Engineering",
			IsActive:   true,
		},
	}

	for _, employee := range employees {
		err := company.AddEmployee(employee)
		if err != nil {
			fmt.Printf("Error adding employee: %s\n", err)
		}
	}

	// Display company information
	fmt.Printf("Company: %s\n", company.Name)
	fmt.Printf("Total Employees: %d\n\n", len(company.Employees))

	// Display departments
	for dept := range company.Departments {
		count, avgSalary, avgService := company.GetDepartmentStats(dept)
		fmt.Printf("Department: %s\n", dept)
		fmt.Printf("  Active Employees: %d\n", count)
		fmt.Printf("  Average Salary: $%.2f\n", avgSalary)
		fmt.Printf("  Average Years of Service: %.1f\n\n", avgService)
	}

	// Perform some operations
	fmt.Println("Operations:")

	// Give a raise
	err := company.UpdateSalary("E001", 100000)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		employee := company.Employees["E001"]
		fmt.Printf("- %s received a raise to $%.2f\n", employee.FullName(), employee.Salary)
	}

	// Transfer an employee
	err = company.TransferEmployee("E002", "Sales")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		employee := company.Employees["E002"]
		fmt.Printf("- %s transferred to %s department\n", employee.FullName(), employee.Department)
	}

	// Terminate an employee
	err = company.TerminateEmployee("E003")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		employee := company.Employees["E003"]
		fmt.Printf("- %s is no longer active\n", employee.FullName())
	}
}
