package service

import (
	"context"
	"shop-api/internal/cache"
	"shop-api/internal/models"
	"shop-api/internal/repository"
)

type ProductService struct {
	repo  *repository.ProductRepository
	cache *cache.RedisCache
}

func NewProductService(repo *repository.ProductRepository, cache *cache.RedisCache) *ProductService {
	return &ProductService{
		repo:  repo,
		cache: cache,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	return s.repo.Create(ctx, req)
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req *models.UpdateProductRequest) error {
	return s.repo.Update(ctx, id, req)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	// Пробуем получить из кэша
	if products, err := s.cache.GetProducts(ctx); err == nil && products != nil {
		return products, nil
	}

	// Если в кэше нет, получаем из БД
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	if err := s.cache.SetProducts(ctx, products); err != nil {
		// Логируем ошибку, но продолжаем работу
		// log.Printf("Failed to cache products: %v", err)
	}

	return products, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, product models.CreateProductRequest) (int, error) {
	id, err := s.repo.Create(ctx, product)
	if err != nil {
		return 0, err
	}

	// Инвалидируем кэш при создании нового продукта
	s.cache.InvalidateProducts(ctx)
	return id, nil
}
