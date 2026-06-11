package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/handlers"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgres://admin:admin123@localhost:5432/interviewdb?sslmode=disable")
	if err != nil {
		t.Fatal("Failed to connect to test DB:", err)
	}
	return db
}

func setupTestRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func TestCreateSession_Success(t *testing.T) {
	// Test ditulis SEBELUM implementasi
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	defer db.Close()
	rdb := setupTestRedis()

	// Insert test user
	var userID int
	db.QueryRow(`INSERT INTO users (email, password) VALUES ('testsession@gmail.com', 'hashed') ON CONFLICT (email) DO UPDATE SET email=EXCLUDED.email RETURNING id`).Scan(&userID)

	handler := handlers.NewSessionHandler(db, rdb)
	r := gin.New()
	r.POST("/sessions", func(c *gin.Context) {
		c.Set("user_id", userID)
		handler.CreateSession(c)
	})

	body := models.CreateSessionRequest{
		Category: "behavioral",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/sessions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var session models.Session
	json.Unmarshal(w.Body.Bytes(), &session)
	assert.Equal(t, "behavioral", session.Category)
	assert.Equal(t, "in_progress", session.Status)

	// Cleanup
	db.Exec("DELETE FROM sessions WHERE user_id = $1", userID)
	db.Exec("DELETE FROM users WHERE id = $1", userID)
}

func TestFinishSession_OptimisticLocking_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	defer db.Close()
	rdb := setupTestRedis()

	var userID int
	db.QueryRow(`INSERT INTO users (email, password) VALUES ('testfinish@gmail.com', 'hashed') ON CONFLICT (email) DO UPDATE SET email=EXCLUDED.email RETURNING id`).Scan(&userID)

	var sessionID int
	db.QueryRow(`INSERT INTO sessions (user_id, category, version) VALUES ($1, 'technical', 1) RETURNING id`, userID).Scan(&sessionID)

	handler := handlers.NewSessionHandler(db, rdb)
	r := gin.New()
	r.PUT("/sessions/:id/finish", func(c *gin.Context) {
		c.Set("user_id", userID)
		handler.FinishSession(c)
	})

	body := map[string]int{"version": 999}
	jsonBody, _ := json.Marshal(body)

	url := fmt.Sprintf("/sessions/%d/finish", sessionID)
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	db.Exec("DELETE FROM sessions WHERE id = $1", sessionID)
	db.Exec("DELETE FROM users WHERE id = $1", userID)
}

func TestSubmitAnswer_SessionNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	defer db.Close()
	rdb := setupTestRedis()

	var userID int
	db.QueryRow(`INSERT INTO users (email, password) VALUES ('testanswer@gmail.com', 'hashed') ON CONFLICT (email) DO UPDATE SET email=EXCLUDED.email RETURNING id`).Scan(&userID)

	handler := handlers.NewSessionHandler(db, rdb)
	r := gin.New()
	r.POST("/sessions/:id/answers", func(c *gin.Context) {
		c.Set("user_id", userID)
		handler.SubmitAnswer(c)
	})

	body := models.SubmitAnswerRequest{
		QuestionID: 1,
		AnswerText: "Test answer",
	}
	jsonBody, _ := json.Marshal(body)

	// Session ID 99999 tidak ada
	req := httptest.NewRequest(http.MethodPost, "/sessions/99999/answers", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Harusnya 404
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Cleanup
	db.Exec("DELETE FROM users WHERE id = $1", userID)
}
