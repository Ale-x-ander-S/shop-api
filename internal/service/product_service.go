package service

import (
	"context"
	"log"
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
	product, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	// Инвалидируем кэш при создании нового продукта
	if err := s.cache.InvalidateProducts(ctx); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req *models.UpdateProductRequest) error {
	err := s.repo.Update(ctx, id, req)
	if err != nil {
		return err
	}

	// Инвалидируем кэш при обновлении продукта
	if err := s.cache.InvalidateProducts(ctx); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Инвалидируем кэш при удалении продукта
	if err := s.cache.InvalidateProducts(ctx); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	// Пробуем получить из кэша
	if products, err := s.cache.GetProducts(ctx); err == nil && products != nil {
		log.Printf("Cache hit: returning %d products from cache", len(products))
		return products, nil
	}

	// Если в кэше нет, получаем из БД
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	if err := s.cache.SetProducts(ctx, products); err != nil {
		log.Printf("Failed to cache products: %v", err)
	} else {
		log.Printf("Cached %d products", len(products))
	}

	return products, nil
}
