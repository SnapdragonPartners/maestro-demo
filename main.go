package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

// homeHandler handles GET requests to the root endpoint
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set Content-Type header to text/html
	w.Header().Set("Content-Type", "text/html")

	// Parse the template file
	tmpl, err := template.ParseFiles("home.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Execute the template and write to response
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// healthHandler handles GET requests to the health endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set Content-Type header to text/plain
	w.Header().Set("Content-Type", "text/plain")

	// Write OK response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// Check if home.html exists
	if _, err := os.Stat("home.html"); os.IsNotExist(err) {
		log.Fatal("home.html template file not found")
	}

	// Register the handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)

	// Start the server
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	log.Println("Server starting on port", port+"...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
