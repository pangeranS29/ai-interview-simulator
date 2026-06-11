package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pangeranS29/ai-interview-simulator/api-service/internal/db"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:admin123@localhost:5432/interviewdb?sslmode=disable"
	}

	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("PostgreSQL not reachable:", err)
	}

	db.Migrate(sqlDB)
	fmt.Println("✅ Migration done!")
}
