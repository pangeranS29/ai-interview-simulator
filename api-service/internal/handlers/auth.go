package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/logger"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

// @Summary Register user baru
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error().Err(err).Msg("Register: invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var user models.User
	err = h.DB.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, created_at",
		req.Email, string(hashed),
	).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		logger.Log.Error().Err(err).Msg("Register: email already exists")
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	logger.Log.Info().Int("user_id", user.ID).Msg("Register: success")
	c.JSON(http.StatusCreated, gin.H{"message": "Register berhasil", "user": user})
}

// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Request"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	var hashedPassword string
	err := h.DB.QueryRow(
		"SELECT id, email, password, created_at FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Email, &hashedPassword, &user.CreatedAt)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret-key-dev"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	logger.Log.Info().Int("user_id", user.ID).Msg("Login: success")
	c.JSON(http.StatusOK, models.LoginResponse{Token: tokenString, User: user})
}

// @Summary Change user password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Change Password Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error().Err(err).Int("user_id", userID).Msg("ChangePassword: invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current password from database
	var user models.User
	var currentHashedPassword string
	err := h.DB.QueryRow(
		"SELECT id, email, password, created_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Email, &currentHashedPassword, &user.CreatedAt)

	if err == sql.ErrNoRows {
		logger.Log.Warn().Int("user_id", userID).Msg("ChangePassword: user not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err != nil {
		logger.Log.Error().Err(err).Int("user_id", userID).Msg("ChangePassword: database error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(currentHashedPassword), []byte(req.OldPassword)); err != nil {
		logger.Log.Warn().Int("user_id", userID).Msg("ChangePassword: wrong old password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password lama tidak sesuai"})
		return
	}

	// Hash new password
	newHashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error().Err(err).Msg("ChangePassword: failed to hash new password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	// Update password in database
	_, err = h.DB.Exec(
		"UPDATE users SET password = $1 WHERE id = $2",
		string(newHashed), userID,
	)

	if err != nil {
		logger.Log.Error().Err(err).Int("user_id", userID).Msg("ChangePassword: failed to update password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	logger.Log.Info().Int("user_id", userID).Msg("ChangePassword: success")
	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}
