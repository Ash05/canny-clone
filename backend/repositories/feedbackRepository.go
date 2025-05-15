package repositories

import (
	"database/sql"
)

type Feedback struct {
	ID          int     `json:"id"`
	BoardID     int     `json:"boardId"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CategoryID  int     `json:"categoryId"`
	Upvotes     int     `json:"upvotes"`
	Downvotes   int     `json:"downvotes"`
	Status      string  `json:"status"`
}

type FeedbackRepository interface {
	GetFeedbacksByBoardID(boardID int) ([]Feedback, error)
	CreateFeedback(feedback *Feedback) error
	UpdateFeedbackVote(id int, isUpvote bool, increment bool) error
	GetFeedbackByID(id int) (*Feedback, error)
	UpdateFeedbackStatus(id int, status string) error
}

type FeedbackRepositoryImpl struct {
	db *sql.DB
}

func NewFeedbackRepository() FeedbackRepository {
	return &FeedbackRepositoryImpl{
		db: GetDB(),
	}
}

func (r *FeedbackRepositoryImpl) GetFeedbacksByBoardID(boardID int) ([]Feedback, error) {
	rows, err := r.db.Query(`
		SELECT id, board_id, title, description, category_id, upvotes, downvotes, COALESCE(status, 'pending')
		FROM feedback 
		WHERE board_id = $1
	`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []Feedback
	for rows.Next() {
		var fb Feedback
		if err := rows.Scan(&fb.ID, &fb.BoardID, &fb.Title, &fb.Description, &fb.CategoryID, &fb.Upvotes, &fb.Downvotes, &fb.Status); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, fb)
	}

	return feedbacks, nil
}

func (r *FeedbackRepositoryImpl) CreateFeedback(feedback *Feedback) error {
	_, err := r.db.Exec(`
		INSERT INTO feedback (board_id, title, description, category_id, upvotes, downvotes, status) 
		VALUES ($1, $2, $3, $4, 0, 0, 'pending')
	`, feedback.BoardID, feedback.Title, feedback.Description, feedback.CategoryID)
	
	return err
}

func (r *FeedbackRepositoryImpl) UpdateFeedbackVote(id int, isUpvote bool, increment bool) error {
	var query string
	if isUpvote {
		if increment {
			query = "UPDATE feedback SET upvotes = upvotes + 1 WHERE id = $1"
		} else {
			query = "UPDATE feedback SET upvotes = upvotes - 1 WHERE id = $1"
		}
	} else {
		if increment {
			query = "UPDATE feedback SET downvotes = downvotes + 1 WHERE id = $1"
		} else {
			query = "UPDATE feedback SET downvotes = downvotes - 1 WHERE id = $1"
		}
	}

	_, err := r.db.Exec(query, id)
	return err
}

func (r *FeedbackRepositoryImpl) GetFeedbackByID(id int) (*Feedback, error) {
	var fb Feedback
	err := r.db.QueryRow(`
		SELECT id, board_id, title, description, category_id, upvotes, downvotes, COALESCE(status, 'pending') 
		FROM feedback 
		WHERE id = $1
	`, id).Scan(&fb.ID, &fb.BoardID, &fb.Title, &fb.Description, &fb.CategoryID, &fb.Upvotes, &fb.Downvotes, &fb.Status)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		return nil, err
	}
	
	return &fb, nil
}

func (r *FeedbackRepositoryImpl) UpdateFeedbackStatus(id int, status string) error {
	_, err := r.db.Exec(`
		UPDATE feedback 
		SET status = $1
		WHERE id = $2
	`, status, id)
	
	return err
}
