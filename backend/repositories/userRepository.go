package repositories

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	Provider  string    `json:"provider"` // e.g., "google", "github"
	Role      string    `json:"role"`     // "app_admin", "stakeholder", "user"
	CreatedAt time.Time `json:"createdAt"`
}

type UserRepository interface {
	FindUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
	GetUserByID(id int) (*User, error)
	UpdateUserRole(userID int, role string) error
	GetUserBoardRoles(userID int) (map[int]string, error)
	AddUserToBoard(userID, boardID int, role string) error
	RemoveUserFromBoard(userID, boardID int) error
	GetBoardMembers(boardID int) ([]*User, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{
		db: GetDB(),
	}
}

func (r *UserRepositoryImpl) FindUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow(`
		SELECT id, email, name, picture, provider, role, created_at 
		FROM users 
		WHERE email = $1
	`, email).Scan(&user.ID, &user.Email, &user.Name, &user.Picture, &user.Provider, &user.Role, &user.CreatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}
	
	return &user, nil
}

func (r *UserRepositoryImpl) CreateUser(user *User) error {
	// Default role is "user" if not specified
	if user.Role == "" {
		user.Role = "user"
	}
	
	_, err := r.db.Exec(`
		INSERT INTO users (email, name, picture, provider, role) 
		VALUES ($1, $2, $3, $4, $5)
	`, user.Email, user.Name, user.Picture, user.Provider, user.Role)
	
	return err
}

func (r *UserRepositoryImpl) GetUserByID(id int) (*User, error) {
	var user User
	err := r.db.QueryRow(`
		SELECT id, email, name, picture, provider, role, created_at 
		FROM users 
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Email, &user.Name, &user.Picture, &user.Provider, &user.Role, &user.CreatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}
	
	return &user, nil
}

// UpdateUserRole updates a user's global role
func (r *UserRepositoryImpl) UpdateUserRole(userID int, role string) error {
	_, err := r.db.Exec(`
		UPDATE users SET role = $1 WHERE id = $2
	`, role, userID)
	return err
}

// GetUserBoardRoles retrieves all boards and the user's role in each
func (r *UserRepositoryImpl) GetUserBoardRoles(userID int) (map[int]string, error) {
	rows, err := r.db.Query(`
		SELECT board_id, role FROM board_members WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boardRoles := make(map[int]string)
	for rows.Next() {
		var boardID int
		var role string
		if err := rows.Scan(&boardID, &role); err != nil {
			return nil, err
		}
		boardRoles[boardID] = role
	}

	return boardRoles, rows.Err()
}

// AddUserToBoard adds a user to a board with a specific role
func (r *UserRepositoryImpl) AddUserToBoard(userID, boardID int, role string) error {
	_, err := r.db.Exec(`
		INSERT INTO board_members (user_id, board_id, role) 
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, board_id) 
		DO UPDATE SET role = $3
	`, userID, boardID, role)
	return err
}

// RemoveUserFromBoard removes a user from a board
func (r *UserRepositoryImpl) RemoveUserFromBoard(userID, boardID int) error {
	_, err := r.db.Exec(`
		DELETE FROM board_members 
		WHERE user_id = $1 AND board_id = $2
	`, userID, boardID)
	return err
}

// GetBoardMembers retrieves all users who are members of a specific board
func (r *UserRepositoryImpl) GetBoardMembers(boardID int) ([]*User, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.email, u.name, u.picture, u.provider, u.role, u.created_at
		FROM users u
		JOIN board_members bm ON u.id = bm.user_id
		WHERE bm.board_id = $1
	`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, 
			&user.Picture, &user.Provider, &user.Role, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, rows.Err()
}
