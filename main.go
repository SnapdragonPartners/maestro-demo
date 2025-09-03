// Package main provides a simple calculator application.
package main

import "fmt"

// Add takes two integers and returns their sum.
func Add(a, b int) int {
	return a + b
}

// main is the entry point of the application.
func main() {
	result := Add(2, 3)
	fmt.Printf("2 + 3 = %d\n", result)
}
