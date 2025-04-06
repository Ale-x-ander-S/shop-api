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
	s.fromCache = false // Сбрасываем флаг в начале метода

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

	return products, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	return s.repo.GetByID(ctx, int(id))
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	return s.GetProductByID(ctx, id)
}

func (s *ProductService) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	// Обновляем кэш новыми данными
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error getting products for cache update: %v", err)
		return product, nil
	}

	if err := s.cache.SetProducts(ctx, products); err != nil {
		log.Printf("Error updating cache: %v", err)
	}

	s.fromCache = false
	return product, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req *models.UpdateProductRequest) error {
	product := &models.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return err
	}

	// Обновляем кэш новыми данными
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error getting products for cache update: %v", err)
		return nil
	}

	if err := s.cache.SetProducts(ctx, products); err != nil {
		log.Printf("Error updating cache: %v", err)
	}

	s.fromCache = false
	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, int(id)); err != nil {
		return err
	}

	// Обновляем кэш новыми данными
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error getting products for cache update: %v", err)
		return nil
	}

	if err := s.cache.SetProducts(ctx, products); err != nil {
		log.Printf("Error updating cache: %v", err)
	}

	s.fromCache = false
	return nil
}
