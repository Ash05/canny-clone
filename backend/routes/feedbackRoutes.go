package routes

import (
	"github.com/gorilla/mux"
	"canny-clone/services"
	"canny-clone/middlewares"
	"net/http"
	"encoding/json"
	"strconv"
)

func RegisterFeedbackRoutes(r *mux.Router) {
	// Public feedback routes with auth
	feedbackRouter := r.PathPrefix("/").Subrouter()
	feedbackRouter.Use(services.AuthMiddleware)
	
	feedbackRouter.HandleFunc("/feedbacks", services.GetFeedbacks).Methods("GET")
	feedbackRouter.HandleFunc("/feedback", services.AddFeedback).Methods("POST")
	feedbackRouter.HandleFunc("/vote", services.VoteFeedback).Methods("POST")
	
	// Stakeholder/admin only routes
	stakeholderRouter := r.PathPrefix("/feedbacks").Subrouter()
	stakeholderRouter.Use(services.AuthMiddleware)
	stakeholderRouter.Use(middlewares.RoleRequired("app_admin", "stakeholder"))
	
	// Update feedback status
	stakeholderRouter.HandleFunc("/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		feedbackID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid feedback ID", http.StatusBadRequest)
			return
		}
		
		var statusUpdate struct {
			Status string `json:"status"` // "pending", "reviewing", "approved", "declined"
		}
		
		if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		
		// Validate status
		validStatuses := map[string]bool{
			"pending":   true,
			"reviewing": true,
			"approved":  true,
			"declined":  true,
		}
		
		if !validStatuses[statusUpdate.Status] {
			http.Error(w, "Invalid status value", http.StatusBadRequest)
			return
		}
		
		// Check if user has permission for this feedback's board
		userID, role, _ := middlewares.GetUserFromRequest(r)
		feedbackRepo := services.GetFeedbackRepository()
		feedback, err := feedbackRepo.GetFeedbackByID(feedbackID)
		
		if err != nil || feedback == nil {
			http.Error(w, "Feedback not found", http.StatusNotFound)
			return
		}
		
		// App admin can update any feedback
		if role != "app_admin" {
			// Check if user is stakeholder for this board
			userRepo := services.GetUserRepository()
			boardRoles, err := userRepo.GetUserBoardRoles(userID)
			
			if err != nil {
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			
			boardRole, exists := boardRoles[feedback.BoardID]
			if !exists || boardRole != "stakeholder" {
				http.Error(w, "Forbidden: Insufficient permissions for this feedback", http.StatusForbidden)
				return
			}
		}
		
		// Update feedback status
		if err := feedbackRepo.UpdateFeedbackStatus(feedbackID, statusUpdate.Status); err != nil {
			http.Error(w, "Failed to update feedback status", http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"id":     strconv.Itoa(feedbackID),
			"status": statusUpdate.Status,
		})
	}).Methods("PUT")
}
