package db

import (
	"database/sql"
	"log"
)

func Migrate(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS questions (
			id SERIAL PRIMARY KEY,
			category VARCHAR(100) NOT NULL,
			content TEXT NOT NULL,
			difficulty VARCHAR(50) DEFAULT 'medium',
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			category VARCHAR(100) NOT NULL,
			status VARCHAR(50) DEFAULT 'in_progress',
			score INTEGER DEFAULT 0,
			version INTEGER DEFAULT 1,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS answers (
			id SERIAL PRIMARY KEY,
			session_id INTEGER REFERENCES sessions(id) ON DELETE CASCADE,
			question_id INTEGER REFERENCES questions(id),
			answer_text TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS feedbacks (
			id SERIAL PRIMARY KEY,
			answer_id INTEGER REFERENCES answers(id) ON DELETE CASCADE,
			score INTEGER,
			strengths TEXT,
			weaknesses TEXT,
			suggestion TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`INSERT INTO questions (category, content, difficulty) VALUES
			('behavioral', 'Ceritakan pengalaman kamu saat menghadapi konflik dalam tim. Bagaimana kamu menyelesaikannya?', 'medium'),
			('behavioral', 'Apa pencapaian terbesar kamu dalam karir atau studi? Mengapa itu penting bagi kamu?', 'easy'),
			('behavioral', 'Bagaimana kamu mengelola deadline yang ketat dengan banyak task sekaligus?', 'medium'),
			('technical', 'Jelaskan perbedaan antara REST API dan GraphQL. Kapan kamu memilih salah satunya?', 'medium'),
			('technical', 'Apa itu Docker dan mengapa containerization penting dalam pengembangan modern?', 'easy'),
			('technical', 'Jelaskan konsep CAP theorem dan bagaimana pengaruhnya terhadap desain sistem terdistribusi.', 'hard'),
			('situational', 'Jika kamu diminta deliver fitur dalam 2 hari tapi estimasi normal adalah 1 minggu, apa yang kamu lakukan?', 'medium'),
			('situational', 'Bagaimana kamu handle situasi di mana requirement berubah di tengah pengerjaan?', 'medium')
		ON CONFLICT DO NOTHING`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatal("Migration failed:", err)
		}
	}
	log.Println("✅ Migration berhasil!")
}
