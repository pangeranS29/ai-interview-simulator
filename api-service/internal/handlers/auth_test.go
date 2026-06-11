package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *AuthHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	handler := NewAuthHandler(db)
	router := gin.New()

	return db, mock, handler, router
}

func TestRegister_Success(t *testing.T) {
	db, mock, handler, router := setupAuthTest(t)
	defer db.Close()

	router.POST("/auth/register", handler.Register)

	// Mock data
	email := "test@example.com"
	password := "password123"
	userID := 1
	createdAt := time.Now()

	// Expected: INSERT query akan dipanggil dan return user
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(email, sqlmock.AnyArg()). // password akan di-hash
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "created_at"}).
				AddRow(userID, email, createdAt),
		)

	// Prepare request
	reqBody := models.RegisterRequest{
		Email:    email,
		Password: password,
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Register berhasil", response["message"])
	assert.NotNil(t, response["user"])

	// Verify all expectations met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_InvalidEmail(t *testing.T) {
	db, _, handler, router := setupAuthTest(t)
	defer db.Close()

	router.POST("/auth/register", handler.Register)

	// Request dengan email invalid
	reqBody := map[string]string{
		"email":    "invalid-email", // bukan format email
		"password": "password123",
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_ShortPassword(t *testing.T) {
	db, _, handler, router := setupAuthTest(t)
	defer db.Close()

	router.POST("/auth/register", handler.Register)

	// Request dengan password terlalu pendek
	reqBody := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "12345", // kurang dari 6 karakter
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	db, mock, handler, router := setupAuthTest(t)
	defer db.Close()

	router.POST("/auth/register", handler.Register)

	email := "existing@example.com"
	password := "password123"

	// Mock: email sudah ada di database
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(email, sqlmock.AnyArg()).
		WillReturnError(sql.ErrNoRows) // Simulate constraint error

	reqBody := models.RegisterRequest{
		Email:    email,
		Password: password,
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "already exists")
}

func TestLogin_Success(t *testing.T) {
	db, mock, handler, router := setupAuthTest(t)
	defer db.Close()

	router.POST("/auth/login", handler.Login)

	email := "test@example.com"
	password := "password123"
	userID := 1
	createdAt := time.Now()

	// Hash password untuk mock
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Mock: SELECT user dari database
	mock.ExpectQuery(`SELECT id, email, password, created_at FROM users WHERE email`).
		WithArgs(email).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "password", "created_at"}).
				AddRow(userID, email, string(hashedPassword), createdAt),
		)

	reqBody := models.LoginRequest{
		Email:    email,
		Password: password,
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, email, response.User.Email)
	assert.Equal(t, userID, response.User.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_UserNotFound(t *testing.T) {
	db, mock, handler, router := setupAuthTest(t)
	defer db.Close()

	router.POST("/auth/login", handler.Login)

	email := "notfound@example.com"
	password := "password123"

	// Mock: user tidak ditemukan
	mock.ExpectQuery(`SELECT id, email, password, created_at FROM users WHERE email`).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	reqBody := models.LoginRequest{
		Email:    email,
		Password: password,
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "salah")
}

func TestLogin_WrongPassword(t *testing.T) {
	db, mock, handler, router := setupAuthTest(t)
	defer db.Close()

	router.POST("/auth/login", handler.Login)

	email := "test@example.com"
	correctPassword := "password123"
	wrongPassword := "wrongpassword"
	userID := 1
	createdAt := time.Now()

	// Hash correct password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	// Mock: user ditemukan tapi password salah
	mock.ExpectQuery(`SELECT id, email, password, created_at FROM users WHERE email`).
		WithArgs(email).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "password", "created_at"}).
				AddRow(userID, email, string(hashedPassword), createdAt),
		)

	reqBody := models.LoginRequest{
		Email:    email,
		Password: wrongPassword, // password salah
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "salah")
}

func TestLogin_InvalidRequest(t *testing.T) {
	db, _, handler, router := setupAuthTest(t)
	defer db.Close()

	router.POST("/auth/login", handler.Login)

	// Request dengan body invalid
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
