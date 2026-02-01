package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type CategoryService struct {
	repository *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repository: repo}
}

func (s *CategoryService) GetCategoriesService() (products []models.Category, err error) {
	return s.repository.GetAllCategories()
}

func (s *CategoryService) AddCategory(newCategory models.Category) (id int, err error) {
	return s.repository.AddCategory(newCategory)
}

func (s *CategoryService) GetCategoryById(id int) (product models.Category, err error) {
	return s.repository.GetCategoryById(id)
}

func (s *CategoryService) EditCategory(newCategory models.Category) (product models.Category, err error) {
	return s.repository.EditCategory(newCategory)
}

func (s *CategoryService) DeleteCategory(id int) (err error) {
	return s.repository.DeleteCategory(id)
}
