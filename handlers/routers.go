package handlers

import (
	"online_store/database"
	"online_store/models"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, storage *database.Storage, productStore *models.ProductStore) {
	RegisterProductRoutes(r, storage, productStore)
	RegisterAuthRoutes(r, storage)
	RegisterOrderRoutes(r, storage, productStore)
}
