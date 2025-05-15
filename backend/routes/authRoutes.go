package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"canny-clone/middlewares"
	"canny-clone/services"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(r *mux.Router) {
	r.HandleFunc("/auth/google/login", services.GetGoogleLoginURL).Methods("GET")
	r.HandleFunc("/auth/google/callback", services.HandleGoogleCallback).Methods("GET")
	
	// Protected routes that require authentication
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.Use(services.AuthMiddleware)
	
	// Profile endpoint
	authRouter.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID, _ := strconv.Atoi(r.Header.Get("User-ID"))
		
		userRepo := services.GetUserRepository()
		user, err := userRepo.GetUserByID(userID)
		
		if err != nil {
			http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
			return
		}
		
		// Get user's board roles
		boardRoles, err := userRepo.GetUserBoardRoles(userID)
		if err != nil {
			http.Error(w, "Failed to fetch user board roles", http.StatusInternalServerError)
			return
		}
		
		// Create profile response with roles
		type ProfileResponse struct {
			ID         int            `json:"id"`
			Email      string         `json:"email"`
			Name       string         `json:"name"`
			Picture    string         `json:"picture,omitempty"`
			Role       string         `json:"role"`
			BoardRoles map[int]string `json:"boardRoles"`
		}
		
		profile := ProfileResponse{
			ID:         user.ID,
			Email:      user.Email,
			Name:       user.Name,
			Picture:    user.Picture,
			Role:       user.Role,
			BoardRoles: boardRoles,
		}
		
		json.NewEncoder(w).Encode(profile)
	}).Methods("GET")
	
	// Admin routes
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middlewares.RoleRequired("app_admin"))
	
	// Update user role - admin only
	adminRouter.HandleFunc("/users/{id}/role", func(w http.ResponseWriter, r *http.Request) {
		var roleRequest struct {
			Role string `json:"role"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&roleRequest); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		
		// Validate role
		if roleRequest.Role != "app_admin" && roleRequest.Role != "stakeholder" && roleRequest.Role != "user" {
			http.Error(w, "Invalid role", http.StatusBadRequest)
			return
		}
		
		// Get user ID from URL
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		
		// Update user role
		userRepo := services.GetUserRepository()
		if err := userRepo.UpdateUserRole(userID, roleRequest.Role); err != nil {
			http.Error(w, "Failed to update user role", http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User role updated successfully"})
	}).Methods("PUT")
}
