package repositories

import (
	"context"
	"errors"
	"kasir-api/models"

	"github.com/jackc/pgx/v5"
)

type CategoryRepository struct {
	db *pgx.Conn
}

func NewCategoryRepository(db *pgx.Conn) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) GetAllCategories() (categories []models.Category, err error) {
	query := "SELECT id, name, description FROM categories"
	rows, err := repo.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (repo *CategoryRepository) AddCategory(newCategoryRequest models.Category) (id int, err error) {
	query := "INSERT INTO categories (name, description) VALUES (@name, @description) RETURNING id"
	args := pgx.NamedArgs{
		"name":        newCategoryRequest.Name,
		"description": newCategoryRequest.Description,
	}

	row := repo.db.QueryRow(context.Background(), query, args)
	err = row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (repo *CategoryRepository) GetCategoryById(id int) (category models.Category, err error) {
	query := "SELECT id, name, description FROM categories WHERE id=$1"
	row := repo.db.QueryRow(context.Background(), query, id)
	err = row.Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return category, errors.New("NOT_FOUND")
		}
		return category, err
	}

	return category, nil
}

func (repo *CategoryRepository) EditCategory(newCategory models.Category) (category models.Category, err error) {
	query := "UPDATE categories SET name=@name, description=@description WHERE id=@id RETURNING id, name, description"
	args := pgx.NamedArgs{
		"name":        newCategory.Name,
		"description": newCategory.Description,
		"id":          newCategory.ID,
	}

	row := repo.db.QueryRow(context.Background(), query, args)

	err = row.Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return category, errors.New("NOT_FOUND")
		}
		return category, err
	}

	return category, nil
}

func (repo *CategoryRepository) DeleteCategory(id int) (err error) {
	query := "DELETE from categories WHERE id=$1"

	commandTag, err := repo.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("NOT_FOUND")
	}

	return nil
}
