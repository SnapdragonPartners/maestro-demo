package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"net/http"
	"os"
	"sync"
	"time"
)

// NumQuestions defines how many questions to select for each quiz session

// hmacSecret is the secret key used for HMAC signing and verification
const hmacSecret = "quiz-app-secret-key-2024"
const NumQuestions = 3

// Question represents a quiz question with multiple-choice answers
type Question struct {
	ID          int      `json:"id"`
	Question    string   `json:"question"`
	Choices     []string `json:"choices"`
	AnswerIndex int      `json:"answer_index"`
	Explanation string   `json:"explanation"`
}

// QuizSession represents an active quiz session
type QuizSession struct {
	ID        string
	Questions []Question
	Current   int
	Score     int
	StartTime time.Time
}

var (
	sessions   = make(map[string]*QuizSession)
	sessionMux sync.RWMutex
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("home.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Message string
	}{
		Message: "Welcome to the Quiz Application!",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// loadQuestions reads and parses the questions.json file
func loadQuestions() ([]Question, error) {
	data, err := os.ReadFile("questions.json")
	if err != nil {
		return nil, err
	}

	var questions []Question
	if err := json.Unmarshal(data, &questions); err != nil {
		return nil, err
	}


	// Validate that correct index is within bounds of answers array
	for i, q := range questions {
		if q.AnswerIndex < 0 || q.AnswerIndex >= len(q.Choices) {
			return nil, fmt.Errorf("question %d (id=%d): correct index %d is out of bounds for answers array of length %d", i, q.ID, q.AnswerIndex, len(q.Choices))
		}
	}
	return questions, nil
}

// selectRandomQuestions randomly selects n questions from the provided slice
func selectRandomQuestions(questions []Question, n int) []Question {
	if len(questions) <= n {
		return questions
	}

	// Create a new random source with current time seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Shuffle and select first n questions
	shuffled := make([]Question, len(questions))
	copy(shuffled, questions)
	
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:n]
}

// quizHandler handles the GET /quiz endpoint to start a new quiz session
func quizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load all questions
	allQuestions, err := loadQuestions()
	if err != nil {
		log.Printf("Error loading questions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Select random questions
	selectedQuestions := selectRandomQuestions(allQuestions, NumQuestions)

	// Create a new session
	sessionID := generateSessionID()
	session := &QuizSession{
		ID:        sessionID,
		Questions: selectedQuestions,
		Current:   0,
		Score:     0,
		StartTime: time.Now(),
	}

	// Store session
	sessionMux.Lock()
	sessions[sessionID] = session
	sessionMux.Unlock()

	// Parse and render template
	tmpl, err := template.ParseFiles("quiz.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Question       Question
		QuestionNumber int
		TotalQuestions int
		Score          int
		SessionID      string
		HMAC           string
	}{
		Question:       selectedQuestions[0],
		QuestionNumber: 1,
		TotalQuestions: len(selectedQuestions),
		Score:          0,
		SessionID:      sessionID,
		HMAC:           generateHMAC(fmt.Sprintf("%s|%d|%d", sessionID, 0, 0)),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// generateSessionID creates a unique session identifier
func generateSessionID() string {
	return time.Now().Format("20060102150405") + "-" + randString(8)
}

// randString generates a random string of specified length
func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	// Check if required files exist
	if _, err := os.Stat("home.html"); os.IsNotExist(err) {
		log.Fatal("home.html not found")
	}

	if _, err := os.Stat("questions.json"); os.IsNotExist(err) {
		log.Fatal("questions.json not found")
	}

	// Register handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/quiz", quizHandler)
	http.HandleFunc("/quiz", quizPostHandler)
	http.HandleFunc("/quiz/results", resultsHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// generateHMAC creates an HMAC-SHA256 signature for the given data
func generateHMAC(data string) string {
	h := hmac.New(sha256.New, []byte(hmacSecret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// verifyHMAC validates an HMAC signature against the given data
func verifyHMAC(data string, signature string) bool {
	expectedMAC := generateHMAC(data)
	return hmac.Equal([]byte(expectedMAC), []byte(signature))
}

// quizPostHandler handles POST /quiz to process answer submissions
func quizPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Extract form fields
	sessionID := r.FormValue("session_id")
	currentStr := r.FormValue("current")
	scoreStr := r.FormValue("score")
	answerStr := r.FormValue("answer")
	hmacSignature := r.FormValue("hmac")

	// Validate required fields
	if sessionID == "" || currentStr == "" || scoreStr == "" || hmacSignature == "" {
		log.Printf("Missing required form fields")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Parse numeric fields
	current, err := strconv.Atoi(currentStr)
	if err != nil {
		log.Printf("Invalid current value: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		log.Printf("Invalid score value: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Verify HMAC signature
	formState := fmt.Sprintf("%s|%d|%d", sessionID, current, score)
	if !verifyHMAC(formState, hmacSignature) {
		log.Printf("HMAC verification failed for session %s", sessionID)
		http.Error(w, "Forbidden: Invalid form signature", http.StatusForbidden)
		return
	}

	// Retrieve session
	sessionMux.RLock()
	session, exists := sessions[sessionID]
	sessionMux.RUnlock()

	if !exists {
		log.Printf("Session not found: %s", sessionID)
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Validate current question index
	if current < 0 || current >= len(session.Questions) {
		log.Printf("Invalid question index: %d", current)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Grade the answer
	currentQuestion := session.Questions[current]
	if answerStr != "" {
		answerIndex, err := strconv.Atoi(answerStr)
		if err == nil && answerIndex == currentQuestion.AnswerIndex {
			score++
		}
	}

	// Update session
	sessionMux.Lock()
	session.Score = score
	session.Current = current + 1
	sessionMux.Unlock()

	// Check if this was the last question
	if current+1 >= len(session.Questions) {
		// Redirect to results page
		http.Redirect(w, r, "/quiz/results?session_id="+url.QueryEscape(sessionID), http.StatusSeeOther)
		return
	}

	// Render next question
	tmpl, err := template.ParseFiles("quiz.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	nextQuestion := session.Questions[current+1]
	data := struct {
		Question       Question
		QuestionNumber int
		TotalQuestions int
		Score          int
		SessionID      string
		HMAC           string
	}{
		Question:       nextQuestion,
		QuestionNumber: current + 2,
		TotalQuestions: len(session.Questions),
		Score:          score,
		SessionID:      sessionID,
		HMAC:           generateHMAC(fmt.Sprintf("%s|%d|%d", sessionID, current+1, score)),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// resultsHandler displays the final quiz results
func resultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session ID from query parameter
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		log.Printf("Missing session_id parameter")
		http.Error(w, "Bad request: missing session_id", http.StatusBadRequest)
		return
	}

	// Retrieve session
	sessionMux.RLock()
	session, exists := sessions[sessionID]
	sessionMux.RUnlock()

	if !exists {
		log.Printf("Session not found: %s", sessionID)
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Parse and render results template
	tmpl, err := template.ParseFiles("results.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Score          int
		TotalQuestions int
		SessionID      string
		Percentage     float64
	}{
		Score:          session.Score,
		TotalQuestions: len(session.Questions),
		SessionID:      sessionID,
		Percentage:     float64(session.Score) / float64(len(session.Questions)) * 100,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
