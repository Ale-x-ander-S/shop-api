package cache

import (
	"context"
	"encoding/json"
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

func (c *RedisCache) GetProducts(ctx context.Context) ([]*models.Product, error) {
	log.Printf("Trying to get products from cache")
	val, err := c.client.Get(ctx, "products").Result()
	if err == redis.Nil {
		log.Printf("Products not found in cache")
		return nil, nil
	}
	if err != nil {
		log.Printf("Error getting products from cache: %v", err)
		return nil, err
	}

	var products []*models.Product
	err = json.Unmarshal([]byte(val), &products)
	if err != nil {
		log.Printf("Error unmarshaling products from cache: %v", err)
		return nil, err
	}
	log.Printf("Successfully got %d products from cache", len(products))
	return products, err
}

func (c *RedisCache) SetProducts(ctx context.Context, products []*models.Product) error {
	log.Printf("Trying to cache %d products", len(products))
	data, err := json.Marshal(products)
	if err != nil {
		log.Printf("Error marshaling products for cache: %v", err)
		return err
	}

	err = c.client.Set(ctx, "products", data, 5*time.Minute).Err()
	if err != nil {
		log.Printf("Error setting products in cache: %v", err)
		return err
	}
	log.Printf("Successfully cached products")
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
