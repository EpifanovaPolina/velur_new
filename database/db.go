package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage struct {
	DB    *sql.DB
	Store *sessions.CookieStore
}

func NewStorage() *Storage {
	godotenv.Load()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgress"),
		getEnv("DB_NAME", "velur"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("БД не отвечает:", err)
	}

	log.Println("Подключение к БД успешно")

	store := sessions.NewCookieStore([]byte(getEnv("SESSION_SECRET", "velur-secret-key-2024")))

	return &Storage{
		DB:    db,
		Store: store,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
