package handlers

import "github.com/gorilla/mux"

func RegisterRoutes(r *mux.Router) {
	RegisterProductRoutes(r)
	RegisterAuthRoutes(r)
	RegisterOrderRoutes(r)
}
