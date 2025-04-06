package cache

import (
	"context"
	"log"
	"shop-api/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	log.Printf("Initializing Redis cache with address: %s", addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // без пароля
		DB:       0,
	})

	// Проверяем подключение
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
	} else {
		log.Printf("Successfully connected to Redis")
	}

	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) GetProducts(ctx context.Context) ([]*models.Product, error) {
	start := time.Now()
	log.Printf("Redis: Trying to get products from cache")

	var products []*models.Product
	err := r.client.Get(ctx, "products").Scan(&products)
	if err != nil {
		log.Printf("Redis: Cache miss for products (took %v)", time.Since(start))
		return nil, err
	}

	log.Printf("Redis: Cache hit for products, found %d items (took %v)", len(products), time.Since(start))
	return products, nil
}

func (r *RedisCache) SetProducts(ctx context.Context, products []*models.Product) error {
	start := time.Now()
	log.Printf("Redis: Setting %d products to cache", len(products))

	err := r.client.Set(ctx, "products", products, 5*time.Minute).Err()
	if err != nil {
		log.Printf("Redis: Error setting products to cache: %v (took %v)", err, time.Since(start))
		return err
	}

	log.Printf("Redis: Successfully cached %d products (took %v)", len(products), time.Since(start))
	return nil
}

func (c *RedisCache) InvalidateProducts(ctx context.Context) error {
	log.Printf("Invalidating products cache")
	err := c.client.Del(ctx, "products").Err()
	if err != nil {
		log.Printf("Error invalidating products cache: %v", err)
		return err
	}
	log.Printf("Successfully invalidated products cache")
	return nil
}
