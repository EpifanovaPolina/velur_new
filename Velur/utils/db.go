package utils

import (
	"database/sql"
	"log"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var DB *sql.DB
var Store = sessions.NewCookieStore([]byte("secret-key-velur-2024"))

func InitDB() {
	var err error
	connStr := "host=localhost port=5433 user=postgres password=1234 dbname=velur sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("БД не отвечает:", err)
	}
	log.Println("Подключение к БД успешно")
}
