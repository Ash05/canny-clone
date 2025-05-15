package services

import (
	"canny-clone/repositories"
	"canny-clone/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetUserBoards returns boards the user has access to based on their role
func GetUserBoards(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from request headers (set by auth middleware)
	userIDStr := r.Header.Get("User-ID")
	userRole := r.Header.Get("User-Role")
	
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	boardRepo := repositories.NewBoardRepository()
	
	var boards []repositories.Board
	
	// Admin sees all boards
	if userRole == "app_admin" {
		boards, err = boardRepo.GetAllBoards()
	} else {
		// Other users only see boards they have access to
		boards, err = boardRepo.GetUserBoards(userID)
	}
	
	if err != nil {
		http.Error(w, "Error fetching boards", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boards)
}

// GetBoard returns a specific board if the user has access to it
func GetBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}
	
	userIDStr := r.Header.Get("User-ID")
	userRole := r.Header.Get("User-Role")
	
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	boardRepo := repositories.NewBoardRepository()
	
	// Check if user has access to this board
	if userRole != "app_admin" {
		userRepo := GetUserRepository()
		boardRoles, err := userRepo.GetUserBoardRoles(userID)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		
		_, hasAccess := boardRoles[boardID]
		if !hasAccess {
			http.Error(w, "Forbidden: No access to this board", http.StatusForbidden)
			return
		}
	}
	
	board, err := boardRepo.GetBoardByID(boardID)
	if err != nil {
		http.Error(w, "Error fetching board", http.StatusInternalServerError)
		return
	}
	
	if board == nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

// CreateBoard creates a new board (admin only)
func CreateBoard(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateBoardName(body.Name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Get user ID (role check already done by middleware)
	userIDStr := r.Header.Get("User-ID")
	userID, _ := strconv.Atoi(userIDStr)

	boardRepo := repositories.NewBoardRepository()
	boardID, err := boardRepo.CreateBoard(body.Name)
	
	if err != nil {
		http.Error(w, "Error creating board", http.StatusInternalServerError)
		return
	}
	
	// Make the admin user a stakeholder of the new board
	userRepo := GetUserRepository()
	if err := userRepo.AddUserToBoard(userID, boardID, "stakeholder"); err != nil {
		http.Error(w, "Error assigning user to board", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": boardID})
}

// UpdateBoard updates a board's details
func UpdateBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}
	
	var body struct {
		Name string `json:"name"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := utils.ValidateBoardName(body.Name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	boardRepo := repositories.NewBoardRepository()
	if err := boardRepo.UpdateBoard(boardID, body.Name); err != nil {
		http.Error(w, "Error updating board", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Board updated successfully"})
}
