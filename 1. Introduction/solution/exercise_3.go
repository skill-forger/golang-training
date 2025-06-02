package main

import (
	"fmt"
	"math"
	"unsafe"
)

func main() {
	// Declare variables of different types
	var (
		intVar    int     = 42
		floatVar  float64 = 3.14159
		boolVar   bool    = true
		stringVar string  = "Hello, Go!"
		runeVar   rune    = 'A'
		byteVar   byte    = 255
	)

	// Display variable values and types
	fmt.Printf("%-10s: %v\t(Type: %T, Size: %d bytes)\n", "Integer", intVar, intVar, unsafe.Sizeof(intVar))
	fmt.Printf("%-10s: %v\t(Type: %T, Size: %d bytes)\n", "Float", floatVar, floatVar, unsafe.Sizeof(floatVar))
	fmt.Printf("%-10s: %v\t(Type: %T, Size: %d bytes)\n", "Boolean", boolVar, boolVar, unsafe.Sizeof(boolVar))
	fmt.Printf("%-10s: %v\t(Type: %T, Size: %d bytes)\n", "String", stringVar, stringVar, unsafe.Sizeof(stringVar))
	fmt.Printf("%-10s: %v (%c)\t(Type: %T, Size: %d bytes)\n", "Rune", runeVar, runeVar, runeVar, unsafe.Sizeof(runeVar))
	fmt.Printf("%-10s: %v\t(Type: %T, Size: %d bytes)\n", "Byte", byteVar, byteVar, unsafe.Sizeof(byteVar))

	// Show limits of different numeric types
	fmt.Println("\n--- Numeric Type Limits ---")
	fmt.Printf("int8    : %d to %d\n", math.MinInt8, math.MaxInt8)
	fmt.Printf("uint8   : 0 to %d\n", math.MaxUint8)
	fmt.Printf("int16   : %d to %d\n", math.MinInt16, math.MaxInt16)
	fmt.Printf("uint16  : 0 to %d\n", math.MaxUint16)
	fmt.Printf("int32   : %d to %d\n", math.MinInt32, math.MaxInt32)
	fmt.Printf("uint32  : 0 to %d\n", math.MaxUint32)
	fmt.Printf("int64   : %d to %d\n", math.MinInt64, math.MaxInt64)
	// Note: MaxUint64 doesn't fit in a signed int64, so we'd need special handling
}
