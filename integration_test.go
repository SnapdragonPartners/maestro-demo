package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
)

func setupTestFiles(t *testing.T) func() {
	homeHTML := `<!DOCTYPE html>
<html><head><title>Quiz App</title></head>
<body><h1>{{.Message}}</h1>
<a href="/quiz">Start Quiz</a>
<a href="/leaderboard">View Leaderboard</a>
</body></html>`
	if err := os.WriteFile("home.html", []byte(homeHTML), 0644); err != nil {
		t.Fatal(err)
	}

	questionsJSON := `[{"id":1,"question":"Q1?","choices":["A","B","C","D"],"answer_index":0,"explanation":"E1"},{"id":2,"question":"Q2?","choices":["A","B","C","D"],"answer_index":1,"explanation":"E2"},{"id":3,"question":"Q3?","choices":["A","B","C","D"],"answer_index":2,"explanation":"E3"},{"id":4,"question":"Q4?","choices":["A","B","C","D"],"answer_index":3,"explanation":"E4"}]`
	if err := os.WriteFile("questions.json", []byte(questionsJSON), 0644); err != nil {
		t.Fatal(err)
	}

	quizHTML := `<!DOCTYPE html>
<html><head><title>Quiz</title></head>
<body>
<div>Question {{.QuestionNumber}} of {{.TotalQuestions}}</div>
<div>Score: {{.Score}}</div>
<div>{{.Question.Question}}</div>
<ul>{{range .Question.Choices}}<li>{{.}}</li>{{end}}</ul>
</body></html>`
	if err := os.WriteFile("quiz.html", []byte(quizHTML), 0644); err != nil {
		t.Fatal(err)
	}

	return func() {
		os.Remove("home.html")
		os.Remove("questions.json")
		os.Remove("quiz.html")
	}
}

func setupTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/quiz", quizHandler)
	return httptest.NewServer(mux)
}

func makeRequest(t *testing.T, url string, method string) (int, string) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	return resp.StatusCode, string(body)
}

func TestIntegrationEndpoints(t *testing.T) {
	cleanup := setupTestFiles(t)
	defer cleanup()
	server := setupTestServer()
	defer server.Close()

	status, body := makeRequest(t, server.URL+"/", "GET")
	if status != http.StatusOK {
		t.Errorf("GET / returned status %d, want %d", status, http.StatusOK)
	}
	if !strings.Contains(body, "Welcome to the Quiz Application!") {
		t.Errorf("GET / body missing welcome message")
	}
	if !strings.Contains(body, "/quiz") {
		t.Errorf("GET / body missing quiz link")
	}
	if !strings.Contains(body, "/leaderboard") {
		t.Errorf("GET / body missing leaderboard link")
	}
}

func TestIntegrationHomeAndHealthWithQuiz(t *testing.T) {
	cleanup := setupTestFiles(t)
	defer cleanup()
	server := setupTestServer()
	defer server.Close()

	status, body := makeRequest(t, server.URL+"/health", "GET")
	if status != http.StatusOK {
		t.Errorf("GET /health returned status %d, want %d", status, http.StatusOK)
	}
	if body != "OK" {
		t.Errorf("GET /health returned body %q, want %q", body, "OK")
	}

	status, body = makeRequest(t, server.URL+"/", "GET")
	if status != http.StatusOK {
		t.Errorf("GET / after health check returned status %d, want %d", status, http.StatusOK)
	}
	if !strings.Contains(body, "Welcome to the Quiz Application!") {
		t.Errorf("GET / after health check missing welcome message")
	}
}

func TestIntegrationQuizWithBaseEndpoints(t *testing.T) {
	cleanup := setupTestFiles(t)
	defer cleanup()
	server := setupTestServer()
	defer server.Close()

	status, _ := makeRequest(t, server.URL+"/quiz", "GET")
	if status != http.StatusOK {
		t.Errorf("GET /quiz returned status %d, want %d", status, http.StatusOK)
	}

	status, body := makeRequest(t, server.URL+"/", "GET")
	if status != http.StatusOK {
		t.Errorf("GET / after quiz returned status %d, want %d", status, http.StatusOK)
	}
	if !strings.Contains(body, "Welcome to the Quiz Application!") {
		t.Errorf("GET / after quiz missing welcome message")
	}

	status, body = makeRequest(t, server.URL+"/health", "GET")
	if status != http.StatusOK {
		t.Errorf("GET /health after quiz returned status %d, want %d", status, http.StatusOK)
	}
	if body != "OK" {
		t.Errorf("GET /health after quiz returned body %q, want %q", body, "OK")
	}
}

func TestIntegrationConcurrentRequests(t *testing.T) {
	cleanup := setupTestFiles(t)
	defer cleanup()
	server := setupTestServer()
	defer server.Close()

	var wg sync.WaitGroup
	errors := make(chan error, 30)

	for i := 0; i < 10; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			status, _ := makeRequest(t, server.URL+"/", "GET")
			if status != http.StatusOK {
				errors <- nil
			}
		}()
		go func() {
			defer wg.Done()
			status, body := makeRequest(t, server.URL+"/health", "GET")
			if status != http.StatusOK || body != "OK" {
				errors <- nil
			}
		}()
		go func() {
			defer wg.Done()
			status, _ := makeRequest(t, server.URL+"/quiz", "GET")
			if status != http.StatusOK {
				errors <- nil
			}
		}()
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		t.Errorf("Concurrent requests had %d errors", len(errors))
	}
}
