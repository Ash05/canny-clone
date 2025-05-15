package repositories

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	FeedbackID int      `json:"feedbackId"`
	UserID    int       `json:"userId"`
	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	CreatedAt time.Time `json:"createdAt"`
	IsLiked   bool      `json:"isLiked,omitempty"`
	IsDisliked bool     `json:"isDisliked,omitempty"`
	Replies   []Reply   `json:"replies,omitempty"`
}

type Reply struct {
	ID        int       `json:"id"`
	CommentID int       `json:"commentId"`
	UserID    int       `json:"userId"`
	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	CreatedAt time.Time `json:"createdAt"`
	IsLiked   bool      `json:"isLiked,omitempty"`
	IsDisliked bool     `json:"isDisliked,omitempty"`
}

type CommentLikeInfo struct {
	ID       int
	IsLike   bool
}

type CommentRepository interface {
	GetCommentsByFeedbackID(feedbackID int, currentUserID int) ([]Comment, error)
	CreateComment(feedbackID int, userID int, content string) (int, error)
	CreateReply(commentID int, userID int, content string) (int, error)
	GetCommentLikeInfo(commentID int, userID int) (*CommentLikeInfo, error)
	GetReplyLikeInfo(replyID int, userID int) (*CommentLikeInfo, error)
	UpsertCommentLike(commentID int, userID int, isLike bool) error
	UpsertReplyLike(replyID int, userID int, isLike bool) error
	DeleteCommentLike(likeID int) error
	UpdateCommentLikeCount(commentID int, isLike bool, increment bool) error
	UpdateReplyLikeCount(replyID int, isLike bool, increment bool) error
}

type CommentRepositoryImpl struct {
	db *sql.DB
}

func NewCommentRepository() CommentRepository {
	return &CommentRepositoryImpl{
		db: GetDB(),
	}
}

func (r *CommentRepositoryImpl) GetCommentsByFeedbackID(feedbackID int, currentUserID int) ([]Comment, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.feedback_id, c.user_id, c.content, c.likes, c.dislikes, c.created_at
		FROM comments c 
		WHERE c.feedback_id = $1
		ORDER BY c.created_at DESC
	`, feedbackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.FeedbackID, &c.UserID, &c.Content, &c.Likes, &c.Dislikes, &c.CreatedAt); err != nil {
			return nil, err
		}
		
		// Get user reactions
		if currentUserID > 0 {
			likeInfo, err := r.GetCommentLikeInfo(c.ID, currentUserID)
			if err == nil && likeInfo != nil {
				c.IsLiked = likeInfo.IsLike
				c.IsDisliked = !likeInfo.IsLike
			}
		}
		
		// Get replies
		replies, err := r.getRepliesForComment(c.ID, currentUserID)
		if err != nil {
			return nil, err
		}
		c.Replies = replies
		
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *CommentRepositoryImpl) getRepliesForComment(commentID int, currentUserID int) ([]Reply, error) {
	rows, err := r.db.Query(`
		SELECT r.id, r.comment_id, r.user_id, r.content, r.likes, r.dislikes, r.created_at
		FROM comment_replies r 
		WHERE r.comment_id = $1
		ORDER BY r.created_at ASC
	`, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []Reply
	for rows.Next() {
		var r Reply
		if err := rows.Scan(&r.ID, &r.CommentID, &r.UserID, &r.Content, &r.Likes, &r.Dislikes, &r.CreatedAt); err != nil {
			return nil, err
		}
		
		// Get user reactions
		if currentUserID > 0 {
			likeInfo, err := r.GetReplyLikeInfo(r.ID, currentUserID)
			if err == nil && likeInfo != nil {
				r.IsLiked = likeInfo.IsLike
				r.IsDisliked = !likeInfo.IsLike
			}
		}
		
		replies = append(replies, r)
	}

	return replies, nil
}

func (r *CommentRepositoryImpl) CreateComment(feedbackID int, userID int, content string) (int, error) {
	var commentID int
	err := r.db.QueryRow(
		"INSERT INTO comments (feedback_id, user_id, content) VALUES ($1, $2, $3) RETURNING id",
		feedbackID, userID, content,
	).Scan(&commentID)
	
	if err != nil {
		return 0, err
	}
	
	return commentID, nil
}

func (r *CommentRepositoryImpl) CreateReply(commentID int, userID int, content string) (int, error) {
	var replyID int
	err := r.db.QueryRow(
		"INSERT INTO comment_replies (comment_id, user_id, content) VALUES ($1, $2, $3) RETURNING id",
		commentID, userID, content,
	).Scan(&replyID)
	
	if err != nil {
		return 0, err
	}
	
	return replyID, nil
}

func (r *CommentRepositoryImpl) GetCommentLikeInfo(commentID int, userID int) (*CommentLikeInfo, error) {
	var info CommentLikeInfo
	err := r.db.QueryRow(
		"SELECT id, is_like FROM comment_likes WHERE comment_id = $1 AND user_id = $2",
		commentID, userID,
	).Scan(&info.ID, &info.IsLike)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No reaction found
		}
		return nil, err
	}
	
	return &info, nil
}

func (r *CommentRepositoryImpl) GetReplyLikeInfo(replyID int, userID int) (*CommentLikeInfo, error) {
	var info CommentLikeInfo
	err := r.db.QueryRow(
		"SELECT id, is_like FROM comment_likes WHERE reply_id = $1 AND user_id = $2",
		replyID, userID,
	).Scan(&info.ID, &info.IsLike)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No reaction found
		}
		return nil, err
	}
	
	return &info, nil
}

func (r *CommentRepositoryImpl) UpsertCommentLike(commentID int, userID int, isLike bool) error {
	_, err := r.db.Exec(
		"INSERT INTO comment_likes (comment_id, user_id, is_like) VALUES ($1, $2, $3) ON CONFLICT (comment_id, user_id) DO UPDATE SET is_like = $3",
		commentID, userID, isLike,
	)
	
	return err
}

func (r *CommentRepositoryImpl) UpsertReplyLike(replyID int, userID int, isLike bool) error {
	_, err := r.db.Exec(
		"INSERT INTO comment_likes (reply_id, user_id, is_like) VALUES ($1, $2, $3) ON CONFLICT (reply_id, user_id) DO UPDATE SET is_like = $3",
		replyID, userID, isLike,
	)
	
	return err
}

func (r *CommentRepositoryImpl) DeleteCommentLike(likeID int) error {
	_, err := r.db.Exec("DELETE FROM comment_likes WHERE id = $1", likeID)
	return err
}

func (r *CommentRepositoryImpl) UpdateCommentLikeCount(commentID int, isLike bool, increment bool) error {
	var query string
	if isLike {
		if increment {
			query = "UPDATE comments SET likes = likes + 1 WHERE id = $1"
		} else {
			query = "UPDATE comments SET likes = likes - 1 WHERE id = $1"
		}
	} else {
		if increment {
			query = "UPDATE comments SET dislikes = dislikes + 1 WHERE id = $1"
		} else {
			query = "UPDATE comments SET dislikes = dislikes - 1 WHERE id = $1"
		}
	}

	_, err := r.db.Exec(query, commentID)
	return err
}

func (r *CommentRepositoryImpl) UpdateReplyLikeCount(replyID int, isLike bool, increment bool) error {
	var query string
	if isLike {
		if increment {
			query = "UPDATE comment_replies SET likes = likes + 1 WHERE id = $1"
		} else {
			query = "UPDATE comment_replies SET likes = likes - 1 WHERE id = $1"
		}
	} else {
		if increment {
			query = "UPDATE comment_replies SET dislikes = dislikes + 1 WHERE id = $1"
		} else {
			query = "UPDATE comment_replies SET dislikes = dislikes - 1 WHERE id = $1"
		}
	}

	_, err := r.db.Exec(query, replyID)
	return err
}
