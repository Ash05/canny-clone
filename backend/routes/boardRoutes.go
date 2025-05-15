package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"canny-clone/middlewares"
	"canny-clone/services"

	"github.com/gorilla/mux"
)

func RegisterBoardRoutes(r *mux.Router) {
	// Authentication required for all board routes
	boardRouter := r.PathPrefix("/").Subrouter()
	boardRouter.Use(services.AuthMiddleware)
	
	// Get boards the user has access to
	boardRouter.HandleFunc("/boards", services.GetUserBoards).Methods("GET")
	
	// Get single board if user has access
	boardRouter.HandleFunc("/boards/{id}", services.GetBoard).Methods("GET")
	
	// Admin and stakeholder routes
	adminBoardRouter := r.PathPrefix("/").Subrouter()
	adminBoardRouter.Use(middlewares.RoleRequired("app_admin", "stakeholder"))
	
	// Create board - only app_admin can create new boards
	adminBoardRouter.HandleFunc("/boards", func(w http.ResponseWriter, r *http.Request) {
		userID, role, err := middlewares.GetUserFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		// Only app_admin can create boards
		if role != "app_admin" {
			http.Error(w, "Forbidden: Only administrators can create boards", http.StatusForbidden)
			return
		}
		
		services.CreateBoard(w, r)
	}).Methods("POST")
	
	// Update board - only app_admin and board stakeholders can update
	adminBoardRouter.HandleFunc("/boards/{id}", func(w http.ResponseWriter, r *http.Request) {
		userID, role, err := middlewares.GetUserFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		vars := mux.Vars(r)
		boardID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid board ID", http.StatusBadRequest)
			return
		}
		
		// App admin can update any board
		if role == "app_admin" {
			services.UpdateBoard(w, r)
			return
		}
		
		// Check if user is a stakeholder for this board
		userRepo := services.GetUserRepository()
		boardRoles, err := userRepo.GetUserBoardRoles(userID)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		
		boardRole, exists := boardRoles[boardID]
		if !exists || boardRole != "stakeholder" {
			http.Error(w, "Forbidden: Insufficient permissions for this board", http.StatusForbidden)
			return
		}
		
		services.UpdateBoard(w, r)
	}).Methods("PUT")
	
	// Board member management - Admin and stakeholders only
	boardMemberRouter := r.PathPrefix("/boards/{id}/members").Subrouter()
	boardMemberRouter.Use(middlewares.RoleRequired("app_admin", "stakeholder"))
	
	// Add member to board
	boardMemberRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		boardID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid board ID", http.StatusBadRequest)
			return
		}
		
		var memberRequest struct {
			Email string `json:"email"`
			Role  string `json:"role"` // "stakeholder" or "user"
		}
		
		if err := json.NewDecoder(r.Body).Decode(&memberRequest); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		
		// Validate role
		if memberRequest.Role != "stakeholder" && memberRequest.Role != "user" {
			http.Error(w, "Invalid role", http.StatusBadRequest)
			return
		}
		
		// Get the user to add
		userRepo := services.GetUserRepository()
		user, err := userRepo.FindUserByEmail(memberRequest.Email)
		if err != nil || user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		
		// Add user to board
		if err := userRepo.AddUserToBoard(user.ID, boardID, memberRequest.Role); err != nil {
			http.Error(w, "Failed to add user to board", http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User added to board successfully"})
	}).Methods("POST")
	
	// Remove member from board
	boardMemberRouter.HandleFunc("/{userID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		boardID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid board ID", http.StatusBadRequest)
			return
		}
		
		userID, err := strconv.Atoi(vars["userID"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		
		// Remove user from board
		userRepo := services.GetUserRepository()
		if err := userRepo.RemoveUserFromBoard(userID, boardID); err != nil {
			http.Error(w, "Failed to remove user from board", http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User removed from board successfully"})
	}).Methods("DELETE")
	
	// List board members
	boardMemberRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		boardID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid board ID", http.StatusBadRequest)
			return
		}
		
		userRepo := services.GetUserRepository()
		members, err := userRepo.GetBoardMembers(boardID)
		if err != nil {
			http.Error(w, "Failed to get board members", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(members)
	}).Methods("GET")
}
