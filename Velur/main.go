package main

import (
	"log"
	"net/http"
	"time"

	"online_store/handlers"
	"online_store/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Инициализация
	utils.InitDB()
	utils.LoadProducts()

	// Создание маршрутов
	r := mux.NewRouter()

	// Статические файлы
	fs := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Регистрация маршрутов
	handlers.RegisterRoutes(r)

	// Настройка сервера
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
