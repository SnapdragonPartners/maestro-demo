package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	
	"os"
	"strconv"
	"sync"
	"time"
)

// NumQuestions defines how many questions to select for each quiz session
const NumQuestions = 3

// hmacSecret is the secret key used for HMAC signing and verification
const hmacSecret = "quiz-app-secret-key-change-in-production"

// generateHMAC generates an HMAC-SHA256 signature for the given data
func generateHMAC(data string) string {
	h := hmac.New(sha256.New, []byte(hmacSecret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
// verifyHMAC verifies that the provided signature matches the HMAC of the data
func verifyHMAC(data, signature string) bool {
	expectedSignature := generateHMAC(data)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

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

	// Generate HMAC signature for form state
	formState := fmt.Sprintf("%s|%d|%d", sessionID, 0, 0)
	hmacSignature := generateHMAC(formState)

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
		CurrentIndex   int
		HMACSignature  string
	}{
		Question:       selectedQuestions[0],
		QuestionNumber: 1,
		TotalQuestions: len(selectedQuestions),
		Score:          0,
		SessionID:      sessionID,
		CurrentIndex:   0,
		HMACSignature:  hmacSignature,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
// quizPostHandler handles POST requests to submit quiz answers
func quizPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	sessionID := r.FormValue("sessionID")
	currentIndexStr := r.FormValue("currentIndex")
	scoreStr := r.FormValue("score")
	hmacSignature := r.FormValue("hmacSignature")
	selectedAnswerStr := r.FormValue("answer")

	if sessionID == "" || currentIndexStr == "" || scoreStr == "" || hmacSignature == "" {
		log.Printf("Missing required form fields")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	currentIndex, err := strconv.Atoi(currentIndexStr)
	if err != nil {
		log.Printf("Invalid currentIndex: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		log.Printf("Invalid score: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	formState := fmt.Sprintf("%s|%d|%d", sessionID, currentIndex, score)
	if !verifyHMAC(formState, hmacSignature) {
		log.Printf("HMAC verification failed for session %s", sessionID)
		http.Error(w, "Invalid request - tampering detected", http.StatusBadRequest)
		return
	}

	log.Printf("HMAC verification successful for session %s", sessionID)

	sessionMux.RLock()
	session, exists := sessions[sessionID]
	sessionMux.RUnlock()

	if !exists {
		log.Printf("Session not found: %s", sessionID)
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	if currentIndex < 0 || currentIndex >= len(session.Questions) {
		log.Printf("Invalid question index: %d", currentIndex)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}


	currentQuestion := session.Questions[currentIndex]
	
	if selectedAnswerStr != "" {
		selectedAnswer, err := strconv.Atoi(selectedAnswerStr)
		if err == nil && selectedAnswer == currentQuestion.AnswerIndex {
			score++
			log.Printf("Correct answer for session %s, question %d. New score: %d", sessionID, currentIndex, score)
		} else {
			log.Printf("Incorrect answer for session %s, question %d. Score remains: %d", sessionID, currentIndex, score)
		}
	}

	sessionMux.Lock()
	session.Score = score
	session.Current = currentIndex + 1
	sessionMux.Unlock()

	log.Printf("Updated session %s: score=%d, current=%d", sessionID, score, currentIndex+1)

	nextIndex := currentIndex + 1
	if nextIndex >= len(session.Questions) {
		log.Printf("Quiz completed for session %s, redirecting to results", sessionID)
		http.Redirect(w, r, "/quiz/results?session="+sessionID, http.StatusSeeOther)
		return
	}

	nextQuestion := session.Questions[nextIndex]
	newFormState := fmt.Sprintf("%s|%d|%d", sessionID, nextIndex, score)
	newHMACSignature := generateHMAC(newFormState)

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
		CurrentIndex   int
		HMACSignature  string
	}{
		Question:       nextQuestion,
		QuestionNumber: nextIndex + 1,
		TotalQuestions: len(session.Questions),
		Score:          score,
		SessionID:      sessionID,
		CurrentIndex:   nextIndex,
		HMACSignature:  newHMACSignature,
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
// resultsHandler displays the final quiz results
func resultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		log.Printf("Missing session parameter")
		http.Error(w, "Bad request - missing session", http.StatusBadRequest)
		return
	}

	sessionMux.RLock()
	session, exists := sessions[sessionID]
	sessionMux.RUnlock()

	if !exists {
		log.Printf("Session not found: %s", sessionID)
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

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
	}{
		Score:          session.Score,
		TotalQuestions: len(session.Questions),
		SessionID:      sessionID,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
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
	http.HandleFunc("/quiz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			quizHandler(w, r)
		} else if r.Method == http.MethodPost {
			quizPostHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
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
