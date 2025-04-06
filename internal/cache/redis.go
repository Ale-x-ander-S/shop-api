package cache

import (
	"context"
	"encoding/json"
	"shop-api/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // без пароля
		DB:       0,
	})

	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) GetProducts(ctx context.Context) ([]models.Product, error) {
	val, err := c.client.Get(ctx, "products").Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var products []models.Product
	err = json.Unmarshal([]byte(val), &products)
	return products, err
}

func (c *RedisCache) SetProducts(ctx context.Context, products []models.Product) error {
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, "products", data, 5*time.Minute).Err()
}

func (c *RedisCache) InvalidateProducts(ctx context.Context) error {
	return c.client.Del(ctx, "products").Err()
}
