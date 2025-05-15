package repositories

import (
	"database/sql"
	"errors"
)

type Board struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type BoardRepository interface {
	GetAllBoards() ([]Board, error)
	GetUserBoards(userID int) ([]Board, error)
	CreateBoard(name string) (int, error)
	GetBoardByID(id int) (*Board, error)
	UpdateBoard(id int, name string) error
}

type BoardRepositoryImpl struct {
	db *sql.DB
}

func NewBoardRepository() BoardRepository {
	return &BoardRepositoryImpl{
		db: GetDB(),
	}
}

func (r *BoardRepositoryImpl) GetAllBoards() ([]Board, error) {
	rows, err := r.db.Query("SELECT id, name FROM boards")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []Board
	for rows.Next() {
		var board Board
		if err := rows.Scan(&board.ID, &board.Name); err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	return boards, nil
}

func (r *BoardRepositoryImpl) CreateBoard(name string) (int, error) {
	var id int
	err := r.db.QueryRow("INSERT INTO boards (name) VALUES ($1) RETURNING id", name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *BoardRepositoryImpl) GetBoardByID(id int) (*Board, error) {
	var board Board
	err := r.db.QueryRow("SELECT id, name FROM boards WHERE id = $1", id).Scan(&board.ID, &board.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("board not found")
		}
		return nil, err
	}
	return &board, nil
}

func (r *BoardRepositoryImpl) GetUserBoards(userID int) ([]Board, error) {
	rows, err := r.db.Query(`
		SELECT b.id, b.name 
		FROM boards b
		JOIN board_members bm ON b.id = bm.board_id
		WHERE bm.user_id = $1
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []Board
	for rows.Next() {
		var board Board
		if err := rows.Scan(&board.ID, &board.Name); err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	return boards, nil
}

func (r *BoardRepositoryImpl) UpdateBoard(id int, name string) error {
	result, err := r.db.Exec("UPDATE boards SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("board not found")
	}

	return nil
}
