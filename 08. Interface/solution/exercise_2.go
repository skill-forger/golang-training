package main

import (
	"fmt"
	"math"
	"sort"
)

// Shape interface defines methods all shapes must implement
type Shape interface {
	Area() float64
	Perimeter() float64
	Name() string
}

// Circle implements the Shape interface
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) Name() string {
	return "Circle"
}

// Rectangle implements the Shape interface
type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) Name() string {
	return "Rectangle"
}

// Triangle implements the Shape interface
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

func (t Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

func (t Triangle) Area() float64 {
	// Heron's formula
	s := t.Perimeter() / 2
	return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}

func (t Triangle) Name() string {
	return "Triangle"
}

// ThreeDimensionalShape extends the Shape interface
type ThreeDimensionalShape interface {
	Shape
	Volume() float64
}

// Sphere implements ThreeDimensionalShape
type Sphere struct {
	Radius float64
}

func (s Sphere) Area() float64 {
	return 4 * math.Pi * s.Radius * s.Radius
}

func (s Sphere) Perimeter() float64 {
	return 2 * math.Pi * s.Radius // Great circle
}

func (s Sphere) Volume() float64 {
	return (4.0 / 3.0) * math.Pi * math.Pow(s.Radius, 3)
}

func (s Sphere) Name() string {
	return "Sphere"
}

// Cube implements ThreeDimensionalShape
type Cube struct {
	Side float64
}

func (c Cube) Area() float64 {
	return 6 * c.Side * c.Side
}

func (c Cube) Perimeter() float64 {
	return 12 * c.Side
}

func (c Cube) Volume() float64 {
	return math.Pow(c.Side, 3)
}

func (c Cube) Name() string {
	return "Cube"
}

// ShapeProcessor provides utility functions for working with shapes
type ShapeProcessor struct{}

// SortByArea sorts shapes by their area
func (sp ShapeProcessor) SortByArea(shapes []Shape) {
	sort.Slice(shapes, func(i, j int) bool {
		return shapes[i].Area() < shapes[j].Area()
	})
}

// PrintShapeInfo displays information about a shape
func (sp ShapeProcessor) PrintShapeInfo(shape Shape) {
	fmt.Printf("%s:\n", shape.Name())
	fmt.Printf("  Area: %.2f\n", shape.Area())
	fmt.Printf("  Perimeter: %.2f\n", shape.Perimeter())

	// Check if it's also a 3D shape
	if threeDShape, ok := shape.(ThreeDimensionalShape); ok {
		fmt.Printf("  Volume: %.2f\n", threeDShape.Volume())
	}
}

// FilterByType returns shapes of a specific type
func (sp ShapeProcessor) FilterByType(shapes []Shape, typeName string) []Shape {
	var result []Shape
	for _, shape := range shapes {
		if shape.Name() == typeName {
			result = append(result, shape)
		}
	}
	return result
}

func main() {
	// Create various shapes
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 6},
		Triangle{SideA: 3, SideB: 4, SideC: 5},
		Sphere{Radius: 3},
		Cube{Side: 4},
	}

	processor := ShapeProcessor{}

	// Print information for each shape
	fmt.Println("All Shapes:")
	for _, shape := range shapes {
		processor.PrintShapeInfo(shape)
		fmt.Println()
	}

	// Sort shapes by area
	processor.SortByArea(shapes)
	fmt.Println("Shapes sorted by area:")
	for _, shape := range shapes {
		fmt.Printf("%s: %.2f\n", shape.Name(), shape.Area())
	}

	// Filter 3D shapes
	var threeDShapes []ThreeDimensionalShape
	for _, shape := range shapes {
		if threeDShape, ok := shape.(ThreeDimensionalShape); ok {
			threeDShapes = append(threeDShapes, threeDShape)
		}
	}

	fmt.Println("\nThree-dimensional shapes:")
	for _, shape := range threeDShapes {
		fmt.Printf("%s - Volume: %.2f\n", shape.Name(), shape.Volume())
	}
}
