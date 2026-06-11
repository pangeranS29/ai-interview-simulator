package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Answer struct {
	ID         int       `json:"id"`
	SessionID  int       `json:"session_id"`
	QuestionID int       `json:"question_id"`
	AnswerText string    `json:"answer_text"`
	CreatedAt  time.Time `json:"created_at"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type FeedbackResult struct {
	Score      int    `json:"score"`
	Strengths  string `json:"strengths"`
	Weaknesses string `json:"weaknesses"`
	Suggestion string `json:"suggestion"`
}

func generateFeedback(answer Answer, questionContent string, apiKey string) (*FeedbackResult, error) {
	client := resty.New()

	prompt := fmt.Sprintf(`Kamu adalah interviewer profesional. Analisis jawaban interview berikut dan berikan feedback dalam format JSON.

Pertanyaan: %s

Jawaban kandidat: %s

Berikan response HANYA dalam format JSON berikut (tanpa markdown, tanpa backticks):
{
  "score": <angka 0-100>,
  "strengths": "<kekuatan jawaban dalam 1-2 kalimat>",
  "weaknesses": "<kelemahan jawaban dalam 1-2 kalimat>",
  "suggestion": "<saran perbaikan spesifik dalam 1-2 kalimat>"
}`, questionContent, answer.AnswerText)

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{Parts: []GeminiPart{{Text: prompt}}},
		},
	}

	var geminiResp GeminiResponse
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&geminiResp).
		Post(fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", apiKey))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("gemini API error: %s", resp.String())
	}

	if len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates returned")
	}

	text := geminiResp.Candidates[0].Content.Parts[0].Text

	var feedback FeedbackResult
	if err := json.Unmarshal([]byte(text), &feedback); err != nil {
		return nil, fmt.Errorf("failed to parse feedback JSON: %v", err)
	}

	return &feedback, nil
}

func main() {
	godotenv.Load()

	// Koneksi PostgreSQL
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:admin123@localhost:5432/interviewdb?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("PostgreSQL not reachable:", err)
	}
	fmt.Println("✅ Worker: PostgreSQL connected!")

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
	fmt.Println("✅ Worker: Redis connected!")

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Fatal("GEMINI_API_KEY tidak ditemukan!")
	}
	fmt.Println("✅ Worker: Gemini API Key loaded!")

	// Subscribe ke channel answer.submitted
	pubsub := rdb.Subscribe(ctx, "answer.submitted")
	defer pubsub.Close()

	fmt.Println("👂 Worker: Listening for answer.submitted events...")

	for msg := range pubsub.Channel() {
		var answer Answer
		if err := json.Unmarshal([]byte(msg.Payload), &answer); err != nil {
			log.Println("Failed to parse answer:", err)
			continue
		}

		fmt.Printf("📨 Event received: Answer ID %d for Session ID %d\n", answer.ID, answer.SessionID)

		// Ambil pertanyaan dari DB
		var questionContent string
		err := db.QueryRow("SELECT content FROM questions WHERE id = $1", answer.QuestionID).Scan(&questionContent)
		if err != nil {
			log.Printf("Failed to get question %d: %v\n", answer.QuestionID, err)
			continue
		}

		// Generate feedback dari Gemini
		feedback, err := generateFeedback(answer, questionContent, geminiKey)
		if err != nil {
			log.Printf("Failed to generate feedback: %v\n", err)
			// Simpan feedback default jika Gemini gagal
			db.Exec(
				`INSERT INTO feedbacks (answer_id, score, strengths, weaknesses, suggestion)
				VALUES ($1, $2, $3, $4, $5)`,
				answer.ID, 50, "Jawaban diterima", "Perlu analisis lebih lanjut", "Coba elaborasi lebih detail",
			)
			continue
		}

		// Simpan feedback ke DB
		_, err = db.Exec(
			`INSERT INTO feedbacks (answer_id, score, strengths, weaknesses, suggestion)
			VALUES ($1, $2, $3, $4, $5)`,
			answer.ID, feedback.Score, feedback.Strengths, feedback.Weaknesses, feedback.Suggestion,
		)
		if err != nil {
			log.Printf("Failed to save feedback: %v\n", err)
			continue
		}

		fmt.Printf("✅ Feedback saved for Answer ID %d — Score: %d\n", answer.ID, feedback.Score)
	}
}
