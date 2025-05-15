package repositories

import (
	"database/sql"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CategoryRepository interface {
	GetAllCategories() ([]Category, error)
}

type CategoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{
		db: GetDB(),
	}
}

func (r *CategoryRepositoryImpl) GetAllCategories() ([]Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}
