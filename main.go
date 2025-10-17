package main

import (
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Health check requested")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/health", healthHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
