package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello from Go container!")
	fmt.Printf("Go version: %s\n", os.Getenv("GOVERSION"))
	fmt.Println("Container build pipeline validation successful.")
}

// Add a simple function for testing
func add(a, b int) int {
	return a + b
}
