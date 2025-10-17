package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	// Ensure home.html exists for testing
	if _, err := os.Stat("home.html"); os.IsNotExist(err) {
		t.Fatal("home.html template file not found")
	}

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	// Call the handler with our request and recorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the Content-Type header
	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, "text/html")
	}

	// Check that the response body contains expected HTML content
	body := rr.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("response body does not contain HTML doctype")
	}

	if !strings.Contains(body, "Welcome to Our Application") {
		t.Error("response body does not contain expected title")
	}

	if !strings.Contains(body, "home.html template") {
		t.Error("response body does not contain expected template content")
	}
}

func TestHomeHandlerMethodNotAllowed(t *testing.T) {
	// Test POST request should return 405 Method Not Allowed
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for POST: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestTemplateExists(t *testing.T) {
	// Verify that the home.html template file exists
	if _, err := os.Stat("home.html"); os.IsNotExist(err) {
		t.Fatal("home.html template file does not exist")
	}
}
