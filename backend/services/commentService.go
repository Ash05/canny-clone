package services

import (
	"canny-clone/repositories"
	"canny-clone/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

// Comment represents a comment on feedback
type Comment struct {
	ID        int    `json:"id"`
	FeedbackID int   `json:"feedbackId"`
	UserID    int    `json:"userId"`
	Content   string `json:"content"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
	CreatedAt string `json:"createdAt"`
	IsLiked   bool   `json:"isLiked,omitempty"`  // Whether current user liked this
	IsDisliked bool  `json:"isDisliked,omitempty"` // Whether current user disliked this
	Replies   []Reply `json:"replies,omitempty"`
}

// Reply represents a reply to a comment
type Reply struct {
	ID        int    `json:"id"`
	CommentID int    `json:"commentId"`
	UserID    int    `json:"userId"`
	Content   string `json:"content"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
	CreatedAt string `json:"createdAt"`
	IsLiked   bool   `json:"isLiked,omitempty"`  
	IsDisliked bool  `json:"isDisliked,omitempty"`
}

// GetComments retrieves all comments for a feedback item
func GetComments(w http.ResponseWriter, r *http.Request) {
	feedbackIDStr := r.URL.Query().Get("feedbackId")
	feedbackID, err := strconv.Atoi(feedbackIDStr)
	if err != nil {
		http.Error(w, "Invalid feedback ID", http.StatusBadRequest)
		return
	}
	
	userID := getUserIDFromRequest(r)

	repo := repositories.NewCommentRepository()
	comments, err := repo.GetCommentsByFeedbackID(feedbackID, userID)
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// Helper to get user ID from request (in a real app, this would come from auth middleware)
func getUserIDFromRequest(r *http.Request) int {
	// This is a placeholder - in a real app, you would get this from your auth system
	// For now, we'll use a dummy user ID for testing
	userIDStr := r.Header.Get("User-ID")
	if userIDStr == "" {
		return 1 // Default test user
	}
	userID, _ := strconv.Atoi(userIDStr)
	return userID
}

// AddComment adds a new comment to a feedback
func AddComment(w http.ResponseWriter, r *http.Request) {
	var body struct {
		FeedbackID int    `json:"feedbackId"`
		Content    string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate comment
	if err := utils.ValidateComment(body.Content); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := getUserIDFromRequest(r)

	repo := repositories.NewCommentRepository()
	commentID, err := repo.CreateComment(body.FeedbackID, userID, body.Content)
	if err != nil {
		http.Error(w, "Error adding comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": commentID})
}

// AddReply adds a reply to a comment
func AddReply(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CommentID int    `json:"commentId"`
		Content   string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate reply
	if err := utils.ValidateComment(body.Content); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := getUserIDFromRequest(r)

	repo := repositories.NewCommentRepository()
	replyID, err := repo.CreateReply(body.CommentID, userID, body.Content)
	if err != nil {
		http.Error(w, "Error adding reply", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": replyID})
}

// LikeComment handles liking/disliking a comment or reply
func LikeComment(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CommentID *int `json:"commentId"`
		ReplyID   *int `json:"replyId"`
		IsLike    bool `json:"isLike"` // true for like, false for dislike
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.CommentID == nil && body.ReplyID == nil {
		http.Error(w, "Either commentId or replyId must be provided", http.StatusBadRequest)
		return
	}

	userID := getUserIDFromRequest(r)
	repo := repositories.NewCommentRepository()
	
	// Transaction handled at repository level
	if body.CommentID != nil {
		// Handle comment like/dislike
		commentID := *body.CommentID
		likeInfo, err := repo.GetCommentLikeInfo(commentID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		
		if likeInfo == nil {
			// New reaction
			if err := repo.UpsertCommentLike(commentID, userID, body.IsLike); err != nil {
				http.Error(w, "Error adding reaction", http.StatusInternalServerError)
				return
			}
			
			if err := repo.UpdateCommentLikeCount(commentID, body.IsLike, true); err != nil {
				http.Error(w, "Error updating counter", http.StatusInternalServerError)
				return
			}
		} else {
			// Existing reaction
			if likeInfo.IsLike == body.IsLike {
				// Same reaction - remove it (toggle off)
				if err := repo.DeleteCommentLike(likeInfo.ID); err != nil {
					http.Error(w, "Error removing reaction", http.StatusInternalServerError)
					return
				}
				
				if err := repo.UpdateCommentLikeCount(commentID, body.IsLike, false); err != nil {
					http.Error(w, "Error updating counter", http.StatusInternalServerError)
					return
				}
			} else {
				// Different reaction - update and adjust counters
				if err := repo.UpsertCommentLike(commentID, userID, body.IsLike); err != nil {
					http.Error(w, "Error updating reaction", http.StatusInternalServerError)
					return
				}
				
				// Decrement old counter
				if err := repo.UpdateCommentLikeCount(commentID, !body.IsLike, false); err != nil {
					http.Error(w, "Error updating counter", http.StatusInternalServerError)
					return
				}
				
				// Increment new counter
				if err := repo.UpdateCommentLikeCount(commentID, body.IsLike, true); err != nil {
					http.Error(w, "Error updating counter", http.StatusInternalServerError)
					return
				}
			}
		}
	} else {
		// Handle reply like/dislike similarly
		replyID := *body.ReplyID
		likeInfo, err := repo.GetReplyLikeInfo(replyID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		
		if likeInfo == nil {
			// New reaction
			if err := repo.UpsertReplyLike(replyID, userID, body.IsLike); err != nil {
				http.Error(w, "Error adding reaction", http.StatusInternalServerError)
				return
			}
			
			if err := repo.UpdateReplyLikeCount(replyID, body.IsLike, true); err != nil {
				http.Error(w, "Error updating counter", http.StatusInternalServerError)
				return
			}
		} else {
			// Existing reaction
			if likeInfo.IsLike == body.IsLike {
				// Same reaction - remove it (toggle off)
				if err := repo.DeleteCommentLike(likeInfo.ID); err != nil {
					http.Error(w, "Error removing reaction", http.StatusInternalServerError)
					return
				}
				
				if err := repo.UpdateReplyLikeCount(replyID, body.IsLike, false); err != nil {
					http.Error(w, "Error updating counter", http.StatusInternalServerError)
					return
				}
			} else {
				// Different reaction - update and adjust counters
				if err := repo.UpsertReplyLike(replyID, userID, body.IsLike); err != nil {
					http.Error(w, "Error updating reaction", http.StatusInternalServerError)
					return
				}
				
				// Decrement old counter
				if err := repo.UpdateReplyLikeCount(replyID, !body.IsLike, false); err != nil {
					http.Error(w, "Error updating counter", http.StatusInternalServerError)
					return
				}
				
				// Increment new counter
				if err := repo.UpdateReplyLikeCount(replyID, body.IsLike, true); err != nil {
					http.Error(w, "Error updating counter", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}
