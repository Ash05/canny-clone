package services

import (
	"canny-clone/repositories"
	"encoding/json"
	"net/http"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	repo := repositories.NewCategoryRepository()
	categories, err := repo.GetAllCategories()
	if err != nil {
		http.Error(w, "Error fetching categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
