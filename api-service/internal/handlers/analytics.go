package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/logger"
)

type AnalyticsHandler struct {
	DB *sql.DB
}

func NewAnalyticsHandler(db *sql.DB) *AnalyticsHandler {
	return &AnalyticsHandler{DB: db}
}

type SessionSummary struct {
	TotalSessions     int     `json:"total_sessions"`
	CompletedSessions int     `json:"completed_sessions"`
	AverageScore      float64 `json:"average_score"`
	BestScore         int     `json:"best_score"`
}

type CategoryStat struct {
	Category      string  `json:"category"`
	TotalSessions int     `json:"total_sessions"`
	AverageScore  float64 `json:"average_score"`
}

type ScoreTrend struct {
	SessionID int    `json:"session_id"`
	Category  string `json:"category"`
	Score     int    `json:"score"`
	CreatedAt string `json:"created_at"`
}

type AnalyticsResponse struct {
	Summary       SessionSummary `json:"summary"`
	CategoryStats []CategoryStat `json:"category_stats"`
	ScoreTrends   []ScoreTrend   `json:"score_trends"`
}

// @Summary Get AI analytics dashboard
// @Tags analytics
// @Produce json
// @Security BearerAuth
// @Success 200 {object} AnalyticsResponse
// @Router /analytics [get]
func (h *AnalyticsHandler) GetAnalytics(c *gin.Context) {
	userID := c.GetInt("user_id")

	// Summary
	var summary SessionSummary
	h.DB.QueryRow(
		`SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'completed') as completed,
			COALESCE(AVG(score) FILTER (WHERE status = 'completed'), 0) as avg_score,
			COALESCE(MAX(score) FILTER (WHERE status = 'completed'), 0) as best_score
		FROM sessions WHERE user_id = $1`,
		userID,
	).Scan(&summary.TotalSessions, &summary.CompletedSessions, &summary.AverageScore, &summary.BestScore)

	// Category stats
	rows, err := h.DB.Query(
		`SELECT category, COUNT(*) as total,
			COALESCE(AVG(score) FILTER (WHERE status = 'completed'), 0) as avg_score
		FROM sessions WHERE user_id = $1
		GROUP BY category ORDER BY avg_score DESC`,
		userID,
	)
	if err != nil {
		logger.Log.Error().Err(err).Msg("GetAnalytics: category stats failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics"})
		return
	}
	defer rows.Close()

	categoryStats := []CategoryStat{}
	for rows.Next() {
		var cs CategoryStat
		rows.Scan(&cs.Category, &cs.TotalSessions, &cs.AverageScore)
		categoryStats = append(categoryStats, cs)
	}

	// Score trends
	trendRows, err := h.DB.Query(
		`SELECT id, category, score, created_at
		FROM sessions WHERE user_id = $1 AND status = 'completed'
		ORDER BY created_at ASC LIMIT 10`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trends"})
		return
	}
	defer trendRows.Close()

	scoreTrends := []ScoreTrend{}
	for trendRows.Next() {
		var st ScoreTrend
		trendRows.Scan(&st.SessionID, &st.Category, &st.Score, &st.CreatedAt)
		scoreTrends = append(scoreTrends, st)
	}

	logger.Log.Info().Int("user_id", userID).Msg("GetAnalytics: success")
	c.JSON(http.StatusOK, AnalyticsResponse{
		Summary:       summary,
		CategoryStats: categoryStats,
		ScoreTrends:   scoreTrends,
	})
}
