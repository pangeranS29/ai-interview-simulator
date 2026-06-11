package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/logger"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/models"
	"github.com/redis/go-redis/v9"
)

type SessionHandler struct {
	DB  *sql.DB
	RDB *redis.Client
}

func NewSessionHandler(db *sql.DB, rdb *redis.Client) *SessionHandler {
	return &SessionHandler{DB: db, RDB: rdb}
}

func (h *SessionHandler) invalidateCache(userID int) {
	cacheKey := fmt.Sprintf("sessions:user:%d", userID)
	h.RDB.Del(context.Background(), cacheKey)
	logger.Log.Info().Int("user_id", userID).Msg("Cache invalidated")
}

// @Summary Buat sesi interview baru
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateSessionRequest true "Create Session Request"
// @Success 201 {object} models.Session
// @Router /sessions [post]
func (h *SessionHandler) CreateSession(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req models.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var session models.Session
	err := h.DB.QueryRow(
		`INSERT INTO sessions (user_id, category)
		VALUES ($1, $2)
		RETURNING id, user_id, category, status, score, version, created_at, updated_at`,
		userID, req.Category,
	).Scan(&session.ID, &session.UserID, &session.Category, &session.Status, &session.Score, &session.Version, &session.CreatedAt, &session.UpdatedAt)

	if err != nil {
		logger.Log.Error().Err(err).Msg("CreateSession: failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	h.invalidateCache(userID)
	logger.Log.Info().Int("session_id", session.ID).Msg("CreateSession: success")
	c.JSON(http.StatusCreated, session)
}

// @Summary Get semua sesi interview
// @Tags sessions
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Session
// @Router /sessions [get]
func (h *SessionHandler) GetSessions(c *gin.Context) {
	userID := c.GetInt("user_id")
	cacheKey := fmt.Sprintf("sessions:user:%d", userID)

	// Cek cache
	cached, err := h.RDB.Get(context.Background(), cacheKey).Result()
	if err == nil {
		logger.Log.Info().Int("user_id", userID).Msg("GetSessions: cache hit")
		c.Header("X-Cache", "HIT")
		c.Data(http.StatusOK, "application/json", []byte(cached))
		return
	}

	rows, err := h.DB.Query(
		`SELECT id, user_id, category, status, score, version, created_at, updated_at
		FROM sessions WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sessions"})
		return
	}
	defer rows.Close()

	sessions := []models.Session{}
	for rows.Next() {
		var s models.Session
		rows.Scan(&s.ID, &s.UserID, &s.Category, &s.Status, &s.Score, &s.Version, &s.CreatedAt, &s.UpdatedAt)
		sessions = append(sessions, s)
	}

	sessionsJSON, _ := json.Marshal(sessions)
	h.RDB.Set(context.Background(), cacheKey, sessionsJSON, 30*time.Second)

	logger.Log.Info().Int("user_id", userID).Int("count", len(sessions)).Msg("GetSessions: cache miss")
	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, sessions)
}

// @Summary Get detail sesi interview
// @Tags sessions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Success 200 {object} models.SessionDetail
// @Router /sessions/{id} [get]
func (h *SessionHandler) GetSessionDetail(c *gin.Context) {
	userID := c.GetInt("user_id")
	sessionID, _ := strconv.Atoi(c.Param("id"))

	var session models.Session
	err := h.DB.QueryRow(
		`SELECT id, user_id, category, status, score, version, created_at, updated_at
		FROM sessions WHERE id = $1 AND user_id = $2`,
		sessionID, userID,
	).Scan(&session.ID, &session.UserID, &session.Category, &session.Status, &session.Score, &session.Version, &session.CreatedAt, &session.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	rows, err := h.DB.Query(
		`SELECT a.id, a.session_id, a.question_id, a.answer_text, a.created_at,
			q.id, q.category, q.content, q.difficulty, q.created_at,
			f.id, f.answer_id, f.score, f.strengths, f.weaknesses, f.suggestion, f.created_at
		FROM answers a
		JOIN questions q ON q.id = a.question_id
		LEFT JOIN feedbacks f ON f.answer_id = a.id
		WHERE a.session_id = $1
		ORDER BY a.created_at ASC`,
		sessionID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get answers"})
		return
	}
	defer rows.Close()

	answers := []models.AnswerWithFeedback{}
	for rows.Next() {
		var awf models.AnswerWithFeedback
		var f models.Feedback
		var fID, fAnswerID, fScore sql.NullInt64
		var fStrengths, fWeaknesses, fSuggestion sql.NullString
		var fCreatedAt sql.NullTime

		rows.Scan(
			&awf.Answer.ID, &awf.Answer.SessionID, &awf.Answer.QuestionID, &awf.Answer.AnswerText, &awf.Answer.CreatedAt,
			&awf.Question.ID, &awf.Question.Category, &awf.Question.Content, &awf.Question.Difficulty, &awf.Question.CreatedAt,
			&fID, &fAnswerID, &fScore, &fStrengths, &fWeaknesses, &fSuggestion, &fCreatedAt,
		)

		if fID.Valid {
			f.ID = int(fID.Int64)
			f.AnswerID = int(fAnswerID.Int64)
			f.Score = int(fScore.Int64)
			f.Strengths = fStrengths.String
			f.Weaknesses = fWeaknesses.String
			f.Suggestion = fSuggestion.String
			f.CreatedAt = fCreatedAt.Time
			awf.Feedback = &f
		}

		answers = append(answers, awf)
	}

	logger.Log.Info().Int("session_id", sessionID).Msg("GetSessionDetail: success")
	c.JSON(http.StatusOK, models.SessionDetail{Session: session, Answers: answers})
}

// @Summary Submit jawaban interview
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Param request body models.SubmitAnswerRequest true "Submit Answer Request"
// @Success 201 {object} models.Answer
// @Router /sessions/{id}/answers [post]
func (h *SessionHandler) SubmitAnswer(c *gin.Context) {
	userID := c.GetInt("user_id")
	sessionID, _ := strconv.Atoi(c.Param("id"))

	var req models.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verifikasi session milik user
	var count int
	h.DB.QueryRow("SELECT COUNT(*) FROM sessions WHERE id = $1 AND user_id = $2", sessionID, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	var answer models.Answer
	err := h.DB.QueryRow(
		`INSERT INTO answers (session_id, question_id, answer_text)
		VALUES ($1, $2, $3)
		RETURNING id, session_id, question_id, answer_text, created_at`,
		sessionID, req.QuestionID, req.AnswerText,
	).Scan(&answer.ID, &answer.SessionID, &answer.QuestionID, &answer.AnswerText, &answer.CreatedAt)

	if err != nil {
		logger.Log.Error().Err(err).Msg("SubmitAnswer: failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit answer"})
		return
	}

	// Publish event ke Redis untuk diproses Worker
	payload, _ := json.Marshal(answer)
	h.RDB.Publish(context.Background(), "answer.submitted", payload)

	h.invalidateCache(userID)
	logger.Log.Info().Int("answer_id", answer.ID).Msg("SubmitAnswer: success, event published")
	c.JSON(http.StatusCreated, answer)
}

// @Summary Selesaikan sesi interview
// @Tags sessions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Param version query int true "Version"
// @Success 200 {object} models.Session
// @Router /sessions/{id}/finish [put]
func (h *SessionHandler) FinishSession(c *gin.Context) {
	userID := c.GetInt("user_id")
	sessionID, _ := strconv.Atoi(c.Param("id"))

	type FinishRequest struct {
		Version int `json:"version" binding:"required"`
	}
	var req FinishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah semua jawaban sudah punya feedback
	var totalAnswers, totalFeedbacks int
	h.DB.QueryRow("SELECT COUNT(*) FROM answers WHERE session_id = $1", sessionID).Scan(&totalAnswers)
	h.DB.QueryRow(
		`SELECT COUNT(*) FROM feedbacks f
		JOIN answers a ON a.id = f.answer_id
		WHERE a.session_id = $1`,
		sessionID,
	).Scan(&totalFeedbacks)

	logger.Log.Info().
		Int("session_id", sessionID).
		Int("total_answers", totalAnswers).
		Int("total_feedbacks", totalFeedbacks).
		Msg("FinishSession: checking feedbacks")

	// Jika feedback belum lengkap, return error dengan info
	if totalAnswers > 0 && totalFeedbacks < totalAnswers {
		logger.Log.Warn().Int("session_id", sessionID).Msg("FinishSession: feedbacks not ready")
		c.JSON(http.StatusTooEarly, gin.H{
			"error":           "AI masih menganalisis jawaban, tunggu beberapa saat",
			"total_answers":   totalAnswers,
			"total_feedbacks": totalFeedbacks,
		})
		return
	}

	// Hitung rata-rata score dari feedbacks
	var avgScore sql.NullFloat64
	h.DB.QueryRow(
		`SELECT AVG(f.score) FROM feedbacks f
		JOIN answers a ON a.id = f.answer_id
		WHERE a.session_id = $1`,
		sessionID,
	).Scan(&avgScore)

	score := 0
	if avgScore.Valid {
		score = int(avgScore.Float64)
	}

	logger.Log.Info().Int("session_id", sessionID).Int("calculated_score", score).Msg("FinishSession: score calculated")

	// Optimistic locking
	var session models.Session
	err := h.DB.QueryRow(
		`UPDATE sessions SET status = 'completed', score = $1, version = version + 1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3 AND version = $4
		RETURNING id, user_id, category, status, score, version, created_at, updated_at`,
		score, sessionID, userID, req.Version,
	).Scan(&session.ID, &session.UserID, &session.Category, &session.Status, &session.Score, &session.Version, &session.CreatedAt, &session.UpdatedAt)

	if err == sql.ErrNoRows {
		logger.Log.Warn().Int("session_id", sessionID).Msg("FinishSession: version conflict")
		c.JSON(http.StatusConflict, gin.H{"error": "Data sudah diubah, silakan refresh"})
		return
	}
	if err != nil {
		logger.Log.Error().Err(err).Msg("FinishSession: failed to update session")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish session"})
		return
	}

	h.invalidateCache(userID)
	logger.Log.Info().Int("session_id", sessionID).Int("score", score).Msg("FinishSession: success")
	c.JSON(http.StatusOK, session)
}
