package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repository *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repository: repo}
}

func (s *ProductService) GetProductsService() (products []models.Product, err error) {
	return s.repository.GetAllProducts()
}

func (s *ProductService) AddProduct(newProduct models.Product) (id int, err error) {
	return s.repository.AddProduct(newProduct)
}

func (s *ProductService) GetProductById(id int) (product models.Product, err error) {
	return s.repository.GetProductById(id)
}

func (s *ProductService) EditProduct(newProduct models.Product) (product models.Product, err error) {
	return s.repository.EditProduct(newProduct)
}

func (s *ProductService) DeleteProduct(id int) (err error) {
	return s.repository.DeleteProduct(id)
}
