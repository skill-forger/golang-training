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
