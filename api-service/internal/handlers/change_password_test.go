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
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// TDD APPROACH: Test ditulis SEBELUM implementasi
// Fitur: Change Password untuk user yang sudah login
// Requirements:
// 1. User harus sudah authenticated
// 2. Harus provide old_password yang benar
// 3. New password harus min 6 karakter
// 4. Berhasil update password di database
// 5. Return success message

func TestChangePassword_Success(t *testing.T) {
	// Arrange: Setup test environment
	db, mock, handler, router := setupAuthTest(t)
	defer db.Close()

	router.PUT("/auth/change-password", func(c *gin.Context) {
		c.Set("user_id", 1) // Simulate authenticated user
		handler.ChangePassword(c)
	})

	userID := 1
	email := "user@example.com"
	oldPassword := "oldpass123"
	newPassword := "newpass456"
	createdAt := time.Now()

	// Hash old password untuk mock
	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

	// Mock: SELECT current password dari database
	mock.ExpectQuery(`SELECT id, email, password, created_at FROM users WHERE id`).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "password", "created_at"}).
				AddRow(userID, email, string(hashedOldPassword), createdAt),
		)

	// Mock: UPDATE password di database
	mock.ExpectExec(`UPDATE users SET password`).
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Act: Send change password request
	reqBody := map[string]string{
		"old_password": oldPassword,
		"new_password": newPassword,
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/auth/change-password", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assert: Should return 200 OK with success message
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"], "berhasil")

	// Verify all database expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangePassword_WrongOldPassword(t *testing.T) {
	// Arrange
	db, mock, handler, router := setupAuthTest(t)
	defer db.Close()

	router.PUT("/auth/change-password", func(c *gin.Context) {
		c.Set("user_id", 1)
		handler.ChangePassword(c)
	})

	userID := 1
	email := "user@example.com"
	correctOldPassword := "oldpass123"
	wrongOldPassword := "wrongpassword"
	newPassword := "newpass456"
	createdAt := time.Now()

	hashedCorrectPassword, _ := bcrypt.GenerateFromPassword([]byte(correctOldPassword), bcrypt.DefaultCost)

	// Mock: SELECT user with correct password hash
	mock.ExpectQuery(`SELECT id, email, password, created_at FROM users WHERE id`).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "password", "created_at"}).
				AddRow(userID, email, string(hashedCorrectPassword), createdAt),
		)

	// Act: Try to change with wrong old password
	reqBody := map[string]string{
		"old_password": wrongOldPassword,
		"new_password": newPassword,
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/auth/change-password", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assert: Should return 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Password lama")
}

func TestChangePassword_NewPasswordTooShort(t *testing.T) {
	// Arrange
	db, _, handler, router := setupAuthTest(t)
	defer db.Close()

	router.PUT("/auth/change-password", func(c *gin.Context) {
		c.Set("user_id", 1)
		handler.ChangePassword(c)
	})

	// Act: Try to change with password less than 6 characters
	reqBody := map[string]string{
		"old_password": "oldpass123",
		"new_password": "12345", // Only 5 characters
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/auth/change-password", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assert: Should return 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChangePassword_MissingOldPassword(t *testing.T) {
	// Arrange
	db, _, handler, router := setupAuthTest(t)
	defer db.Close()

	router.PUT("/auth/change-password", func(c *gin.Context) {
		c.Set("user_id", 1)
		handler.ChangePassword(c)
	})

	// Act: Send request without old_password
	reqBody := map[string]string{
		"new_password": "newpass456",
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/auth/change-password", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assert: Should return 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChangePassword_UserNotFound(t *testing.T) {
	// Arrange
	db, mock, handler, router := setupAuthTest(t)
	defer db.Close()

	router.PUT("/auth/change-password", func(c *gin.Context) {
		c.Set("user_id", 999) // User ID yang tidak ada
		handler.ChangePassword(c)
	})

	// Mock: User not found in database
	mock.ExpectQuery(`SELECT id, email, password, created_at FROM users WHERE id`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	// Act
	reqBody := map[string]string{
		"old_password": "oldpass123",
		"new_password": "newpass456",
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/auth/change-password", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assert: Should return 404 Not Found
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestChangePassword_DatabaseError(t *testing.T) {
	// Arrange
	db, mock, handler, router := setupAuthTest(t)
	defer db.Close()

	router.PUT("/auth/change-password", func(c *gin.Context) {
		c.Set("user_id", 1)
		handler.ChangePassword(c)
	})

	userID := 1
	email := "user@example.com"
	oldPassword := "oldpass123"
	newPassword := "newpass456"
	createdAt := time.Now()

	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

	// Mock: SELECT success
	mock.ExpectQuery(`SELECT id, email, password, created_at FROM users WHERE id`).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "password", "created_at"}).
				AddRow(userID, email, string(hashedOldPassword), createdAt),
		)

	// Mock: UPDATE fails with database error
	mock.ExpectExec(`UPDATE users SET password`).
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnError(sql.ErrConnDone)

	// Act
	reqBody := map[string]string{
		"old_password": oldPassword,
		"new_password": newPassword,
	}
	bodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/auth/change-password", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assert: Should return 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
