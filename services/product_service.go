package services

import "kasir-api/repositories"

type ProductService struct {
	repository *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repository: repo}
}
