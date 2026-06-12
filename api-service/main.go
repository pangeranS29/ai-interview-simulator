// @title           AI-Powered Interview Simulator API
// @version         1.0
// @description     API untuk AI-Powered Interview Simulator
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "github.com/pangeranS29/ai-interview-simulator/api-service/docs"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/db"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/handlers"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/logger"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/middleware"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	logger.Init()
	logger.Log.Info().Msg("Starting AI Interview Simulator API...")

	// Koneksi PostgreSQL
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:admin123@localhost:5432/interviewdb?sslmode=disable"
	}

	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("PostgreSQL not reachable:", err)
	}
	logger.Log.Info().Msg("✅ PostgreSQL connected!")

	db.Migrate(sqlDB)

	// Koneksi Redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	var rdb *redis.Client
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		rdb = redis.NewClient(&redis.Options{Addr: redisURL})
	} else {
		rdb = redis.NewClient(opt)
	}

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis not reachable:", err)
	}
	logger.Log.Info().Msg("✅ Redis connected!")

	// Setup Gin
	r := gin.Default()

	// CORS
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:3002",
	}
	if prodOrigins := os.Getenv("CORS_ORIGINS"); prodOrigins != "" {
		allowedOrigins = append(allowedOrigins, prodOrigins)
	}
	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-Cache"},
		AllowCredentials: true,
	}))

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes
	authHandler := handlers.NewAuthHandler(sqlDB)
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		// Protected auth routes
		authProtected := auth.Group("/")
		authProtected.Use(middleware.AuthMiddleware())
		{
			authProtected.PUT("/change-password", authHandler.ChangePassword)
		}
	}

	// Protected routes
	sessionHandler := handlers.NewSessionHandler(sqlDB, rdb)
	questionHandler := handlers.NewQuestionHandler(sqlDB)
	analyticsHandler := handlers.NewAnalyticsHandler(sqlDB)

	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		// Questions
		api.GET("/questions", questionHandler.GetQuestions)

		// Sessions
		api.POST("/sessions", sessionHandler.CreateSession)
		api.GET("/sessions", sessionHandler.GetSessions)
		api.GET("/sessions/:id", sessionHandler.GetSessionDetail)
		api.POST("/sessions/:id/answers", sessionHandler.SubmitAnswer)
		api.PUT("/sessions/:id/finish", sessionHandler.FinishSession)

		// Analytics
		api.GET("/analytics", analyticsHandler.GetAnalytics)
	}

	logger.Log.Info().Str("port", "8080").Msg("🚀 Server running!")
	fmt.Println("📖 Swagger UI: http://localhost:8080/swagger/index.html")
	r.Run(":8080")
}
