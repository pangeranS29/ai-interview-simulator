package models

import "time"

type Question struct {
	ID         int       `json:"id"`
	Category   string    `json:"category"`
	Content    string    `json:"content"`
	Difficulty string    `json:"difficulty"`
	CreatedAt  time.Time `json:"created_at"`
}

type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Category  string    `json:"category"`
	Status    string    `json:"status"`
	Score     int       `json:"score"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSessionRequest struct {
	Category string `json:"category" binding:"required"`
}

type Answer struct {
	ID         int       `json:"id"`
	SessionID  int       `json:"session_id"`
	QuestionID int       `json:"question_id"`
	AnswerText string    `json:"answer_text"`
	CreatedAt  time.Time `json:"created_at"`
}

type SubmitAnswerRequest struct {
	QuestionID int    `json:"question_id" binding:"required"`
	AnswerText string `json:"answer_text" binding:"required"`
}

type Feedback struct {
	ID         int       `json:"id"`
	AnswerID   int       `json:"answer_id"`
	Score      int       `json:"score"`
	Strengths  string    `json:"strengths"`
	Weaknesses string    `json:"weaknesses"`
	Suggestion string    `json:"suggestion"`
	CreatedAt  time.Time `json:"created_at"`
}

type AnswerWithFeedback struct {
	Answer   Answer    `json:"answer"`
	Question Question  `json:"question"`
	Feedback *Feedback `json:"feedback"`
}

type SessionDetail struct {
	Session Session              `json:"session"`
	Answers []AnswerWithFeedback `json:"answers"`
}
