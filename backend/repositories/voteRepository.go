package repositories

import (
	"database/sql"
)

type Vote struct {
	ID         int    `json:"id"`
	FeedbackID int    `json:"feedbackId"`
	UserID     int    `json:"userId"`
	VoteType   string `json:"voteType"` // 'upvote' or 'downvote'
}

type VoteRepository interface {
	GetVoteByFeedbackAndUser(feedbackID, userID int) (*Vote, error)
	CreateVote(vote *Vote) error
	UpdateVote(id int, voteType string) error
	DeleteVote(id int) error
}

type VoteRepositoryImpl struct {
	db *sql.DB
}

func NewVoteRepository() VoteRepository {
	return &VoteRepositoryImpl{
		db: GetDB(),
	}
}

func (r *VoteRepositoryImpl) GetVoteByFeedbackAndUser(feedbackID, userID int) (*Vote, error) {
	var vote Vote
	err := r.db.QueryRow(
		"SELECT id, feedback_id, user_id, vote_type FROM votes WHERE feedback_id = $1 AND user_id = $2",
		feedbackID, userID,
	).Scan(&vote.ID, &vote.FeedbackID, &vote.UserID, &vote.VoteType)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No vote found
		}
		return nil, err
	}
	
	return &vote, nil
}

func (r *VoteRepositoryImpl) CreateVote(vote *Vote) error {
	_, err := r.db.Exec(
		"INSERT INTO votes (feedback_id, user_id, vote_type) VALUES ($1, $2, $3)",
		vote.FeedbackID, vote.UserID, vote.VoteType,
	)
	return err
}

func (r *VoteRepositoryImpl) UpdateVote(id int, voteType string) error {
	_, err := r.db.Exec(
		"UPDATE votes SET vote_type = $1 WHERE id = $2",
		voteType, id,
	)
	return err
}

func (r *VoteRepositoryImpl) DeleteVote(id int) error {
	_, err := r.db.Exec("DELETE FROM votes WHERE id = $1", id)
	return err
}
