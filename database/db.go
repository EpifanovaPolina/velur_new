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
		getEnv("DB_HOST", ${DB_HOST}),
		getEnv("DB_PORT", ${DB_PORT}),
		getEnv("DB_USER", ${DB_USER}),
		getEnv("DB_PASSWORD", ${DB_PASSWORD}),
		getEnv("DB_NAME", ${DB_NAME}))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("БД не отвечает:", err)
	}

	log.Println("Подключение к БД успешно")

	store := sessions.NewCookieStore([]byte(getEnv("SESSION_SECRET", ${SESSION_SECRET})))

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
