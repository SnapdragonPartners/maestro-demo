package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	// Create a test home.html file
	homeHTML := `<!DOCTYPE html><html><body><h1>{{.Message}}</h1></body></html>`
	if err := os.WriteFile("home.html", []byte(homeHTML), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("home.html")

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(), "Welcome to the Quiz Application!") {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestQuizHandler(t *testing.T) {
	// Create test questions.json file
	questionsJSON := `[
		{
			"id": 1,
			"question": "Test question 1?",
			"answers": ["A", "B", "C", "D"],
			"correct": 0
		},
		{
			"id": 2,
			"question": "Test question 2?",
			"answers": ["A", "B", "C", "D"],
			"correct": 1
		},
		{
			"id": 3,
			"question": "Test question 3?",
			"answers": ["A", "B", "C", "D"],
			"correct": 2
		}
	]`
	if err := os.WriteFile("questions.json", []byte(questionsJSON), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("questions.json")

	// Create test quiz.html file
	quizHTML := `<!DOCTYPE html><html><body><h1>{{.Question.Question}}</h1><p>Score: {{.Score}}</p></body></html>`
	if err := os.WriteFile("quiz.html", []byte(quizHTML), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("quiz.html")

	req, err := http.NewRequest("GET", "/quiz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(quizHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("handler returned wrong content type: got %v want text/html", contentType)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Score: 0") {
		t.Errorf("handler should display initial score of 0")
	}

	if !strings.Contains(body, "Test question") {
		t.Errorf("handler should display a test question")
	}
}

func TestQuizHandlerMethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("POST", "/quiz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(quizHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestLoadQuestions(t *testing.T) {
	// Create test questions.json file
	questionsJSON := `[
		{
			"id": 1,
			"question": "Test question?",
			"answers": ["A", "B", "C", "D"],
			"correct": 0
		}
	]`
	if err := os.WriteFile("questions.json", []byte(questionsJSON), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("questions.json")

	questions, err := loadQuestions()
	if err != nil {
		t.Errorf("loadQuestions() returned error: %v", err)
	}

	if len(questions) != 1 {
		t.Errorf("loadQuestions() returned wrong number of questions: got %v want 1", len(questions))
	}

	if questions[0].Question != "Test question?" {
		t.Errorf("loadQuestions() returned wrong question text: got %v", questions[0].Question)
	}
}

func TestSelectRandomQuestions(t *testing.T) {
	questions := []Question{
		{ID: 1, Question: "Q1", Answers: []string{"A", "B", "C", "D"}, Correct: 0},
		{ID: 2, Question: "Q2", Answers: []string{"A", "B", "C", "D"}, Correct: 1},
		{ID: 3, Question: "Q3", Answers: []string{"A", "B", "C", "D"}, Correct: 2},
		{ID: 4, Question: "Q4", Answers: []string{"A", "B", "C", "D"}, Correct: 3},
		{ID: 5, Question: "Q5", Answers: []string{"A", "B", "C", "D"}, Correct: 0},
	}

	selected := selectRandomQuestions(questions, 3)

	if len(selected) != 3 {
		t.Errorf("selectRandomQuestions() returned wrong number of questions: got %v want 3", len(selected))
	}

	// Verify all selected questions are from the original set
	for _, q := range selected {
		found := false
		for _, orig := range questions {
			if q.ID == orig.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("selectRandomQuestions() returned question not in original set: %v", q.ID)
		}
	}
}

func TestSelectRandomQuestionsFewerThanRequested(t *testing.T) {
	questions := []Question{
		{ID: 1, Question: "Q1", Answers: []string{"A", "B", "C", "D"}, Correct: 0},
		{ID: 2, Question: "Q2", Answers: []string{"A", "B", "C", "D"}, Correct: 1},
	}

	selected := selectRandomQuestions(questions, 5)

	if len(selected) != 2 {
		t.Errorf("selectRandomQuestions() should return all questions when n > len: got %v want 2", len(selected))
	}
}

func TestHomeHandlerMethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestHealthHandlerMethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("POST", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}
