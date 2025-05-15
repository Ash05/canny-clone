package routes

import (
	"github.com/gorilla/mux"
	"canny-clone/services"
	"net/http"
)

func RegisterCategoryRoutes(r *mux.Router) {
	r.HandleFunc("/categories", services.GetCategories).Methods("GET")
}
