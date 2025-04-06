package service

import (
	"context"
	"shop-api/internal/models"
	"shop-api/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *models.CreateProductRequest) (int64, error) {
	return s.repo.Create(ctx, product)
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, product *models.UpdateProductRequest) error {
	return s.repo.Update(ctx, id, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	return s.repo.GetAll(ctx)
}
