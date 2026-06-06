package main

import (
	"log"
	"net/http"
	"time"

	"online_store/database"
	"online_store/handlers"
	"online_store/models"

	"github.com/gorilla/mux"
)

func main() {
	storage := database.NewStorage()
	productStore := models.NewProductStore()
	productStore.LoadProducts(storage.DB)

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	handlers.RegisterRoutes(r, storage, productStore)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
