package services

import (
	"canny-clone/repositories"
	"canny-clone/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	boardIDStr := r.URL.Query().Get("boardId")
	boardID, err := strconv.Atoi(boardIDStr)
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}

	repo := repositories.NewFeedbackRepository()
	feedbacks, err := repo.GetFeedbacksByBoardID(boardID)
	if err != nil {
		http.Error(w, "Error fetching feedbacks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedbacks)
}

func AddFeedback(w http.ResponseWriter, r *http.Request) {
	var body struct {
		BoardID     int    `json:"boardId"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CategoryID  int    `json:"categoryId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateFeedback(body.Title, body.Description, body.CategoryID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	feedback := &repositories.Feedback{
		BoardID:     body.BoardID,
		Title:       body.Title,
		Description: body.Description,
		CategoryID:  body.CategoryID,
	}

	repo := repositories.NewFeedbackRepository()
	if err := repo.CreateFeedback(feedback); err != nil {
		http.Error(w, "Error adding feedback", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Helper to get user ID from request (in a real app, this would come from auth middleware)
func getUserIDFromRequest(r *http.Request) int {
	userIDStr := r.Header.Get("User-ID")
	if userIDStr == "" {
		return 0
	}
	userID, _ := strconv.Atoi(userIDStr)
	return userID
}

// Get a reference to the feedback repository
func GetFeedbackRepository() repositories.FeedbackRepository {
	return repositories.NewFeedbackRepository()
}

func VoteFeedback(w http.ResponseWriter, r *http.Request) {
	var body struct {
		FeedbackID int    `json:"feedbackId"`
		VoteType   string `json:"voteType"` // "upvote" or "downvote"
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate vote type
	if body.VoteType != "upvote" && body.VoteType != "downvote" {
		http.Error(w, "Invalid vote type", http.StatusBadRequest)
		return
	}

	// Get user ID from request (in a real app, this would come from auth middleware)
	userID := getUserIDFromRequest(r)

	// Initialize repositories
	voteRepo := repositories.NewVoteRepository()
	feedbackRepo := repositories.NewFeedbackRepository()

	// Check if user has already voted on this feedback
	existingVote, err := voteRepo.GetVoteByFeedbackAndUser(body.FeedbackID, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Handle voting logic
	if existingVote == nil {
		// User hasn't voted before, create new vote
		newVote := &repositories.Vote{
			FeedbackID: body.FeedbackID,
			UserID:     userID,
			VoteType:   body.VoteType,
		}
		
		if err := voteRepo.CreateVote(newVote); err != nil {
			http.Error(w, "Error creating vote", http.StatusInternalServerError)
			return
		}
		
		// Update feedback vote count
		if err := feedbackRepo.UpdateFeedbackVote(body.FeedbackID, body.VoteType == "upvote", true); err != nil {
			http.Error(w, "Error updating vote count", http.StatusInternalServerError)
			return
		}
	} else if existingVote.VoteType == body.VoteType {
		// User is clicking the same vote type again - toggle off (remove vote)
		if err := voteRepo.DeleteVote(existingVote.ID); err != nil {
			http.Error(w, "Error removing vote", http.StatusInternalServerError)
			return
		}
		
		// Decrement vote count
		if err := feedbackRepo.UpdateFeedbackVote(body.FeedbackID, body.VoteType == "upvote", false); err != nil {
			http.Error(w, "Error updating vote count", http.StatusInternalServerError)
			return
		}
	} else {
		// User is changing their vote from upvote to downvote or vice versa
		if err := voteRepo.UpdateVote(existingVote.ID, body.VoteType); err != nil {
			http.Error(w, "Error updating vote", http.StatusInternalServerError)
			return
		}
		
		// Decrement old vote type count
		if err := feedbackRepo.UpdateFeedbackVote(body.FeedbackID, existingVote.VoteType == "upvote", false); err != nil {
			http.Error(w, "Error updating vote count", http.StatusInternalServerError)
			return
		}
		
		// Increment new vote type count
		if err := feedbackRepo.UpdateFeedbackVote(body.FeedbackID, body.VoteType == "upvote", true); err != nil {
			http.Error(w, "Error updating vote count", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
