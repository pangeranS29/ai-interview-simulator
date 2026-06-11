package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/logger"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/models"
)

type QuestionHandler struct {
	DB *sql.DB
}

func NewQuestionHandler(db *sql.DB) *QuestionHandler {
	return &QuestionHandler{DB: db}
}

// @Summary Get pertanyaan berdasarkan kategori
// @Tags questions
// @Produce json
// @Security BearerAuth
// @Param category query string false "Category (behavioral/technical/situational)"
// @Success 200 {array} models.Question
// @Router /questions [get]
func (h *QuestionHandler) GetQuestions(c *gin.Context) {
	category := c.Query("category")

	var rows *sql.Rows
	var err error

	if category != "" {
		rows, err = h.DB.Query(
			`SELECT id, category, content, difficulty, created_at
			FROM questions WHERE category = $1 ORDER BY RANDOM() LIMIT 5`,
			category,
		)
	} else {
		rows, err = h.DB.Query(
			`SELECT id, category, content, difficulty, created_at
			FROM questions ORDER BY RANDOM() LIMIT 5`,
		)
	}

	if err != nil {
		logger.Log.Error().Err(err).Msg("GetQuestions: failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get questions"})
		return
	}
	defer rows.Close()

	questions := []models.Question{}
	for rows.Next() {
		var q models.Question
		rows.Scan(&q.ID, &q.Category, &q.Content, &q.Difficulty, &q.CreatedAt)
		questions = append(questions, q)
	}

	logger.Log.Info().Str("category", category).Int("count", len(questions)).Msg("GetQuestions: success")
	c.JSON(http.StatusOK, questions)
}
