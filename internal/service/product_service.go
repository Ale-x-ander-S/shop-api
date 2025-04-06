package service

import (
	"context"
	"log"
	"shop-api/internal/cache"
	"shop-api/internal/models"
	"shop-api/internal/repository"
)

type ProductService struct {
	repo      repository.ProductRepository
	cache     *cache.RedisCache
	fromCache bool
}

func NewProductService(repo repository.ProductRepository, cache *cache.RedisCache) *ProductService {
	return &ProductService{
		repo:      repo,
		cache:     cache,
		fromCache: false,
	}
}

func (s *ProductService) IsFromCache() bool {
	return s.fromCache
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	// Пробуем получить из кэша
	products, err := s.cache.GetProducts(ctx)
	if err == nil && len(products) > 0 {
		log.Printf("Cache hit: returning %d products from cache", len(products))
		s.fromCache = true
		return products, nil
	}

	// Если в кэше нет, получаем из БД
	products, err = s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error getting products from database: %v", err)
		return nil, err
	}

	// Сохраняем в кэш
	if err := s.cache.SetProducts(ctx, products); err != nil {
		log.Printf("Error caching products: %v", err)
	}

	s.fromCache = false
	return products, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	createdProduct, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	// Обновляем кэш новыми данными
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error getting products for cache update: %v", err)
		return createdProduct, nil // Возвращаем созданный продукт даже если не удалось обновить кэш
	}

	if err := s.cache.SetProducts(ctx, products); err != nil {
		log.Printf("Error updating cache: %v", err)
	}

	s.fromCache = false
	return createdProduct, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req *models.UpdateProductRequest) error {
	err := s.repo.Update(ctx, id, req)
	if err != nil {
		return err
	}

	// Обновляем кэш новыми данными
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error getting products for cache update: %v", err)
		return nil // Возвращаем nil даже если не удалось обновить кэш
	}

	if err := s.cache.SetProducts(ctx, products); err != nil {
		log.Printf("Error updating cache: %v", err)
	}

	s.fromCache = false
	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Обновляем кэш новыми данными
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error getting products for cache update: %v", err)
		return nil // Возвращаем nil даже если не удалось обновить кэш
	}

	if err := s.cache.SetProducts(ctx, products); err != nil {
		log.Printf("Error updating cache: %v", err)
	}

	s.fromCache = false
	return nil
}
