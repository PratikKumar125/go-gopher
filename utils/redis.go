package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache() *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", 
		DB:       0,
	})
	fmt.Println("REDIS INITIALIZED")
	return &Cache{client: client}
}

func (rd *Cache) Client() *redis.Client {
	return rd.client
}

func (rd *Cache) AddKeyWithTTL(ctx context.Context, key string, value interface{}, expiry string) error {
	duration, err := time.ParseDuration(expiry)
	if err != nil {
		return fmt.Errorf("invalid expiry duration: %v", err)
	}

	err = rd.client.Set(ctx, key, value, duration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key in Redis: %v", err)
	}

	return nil
}

func (rd *Cache) RemoveKey(ctx context.Context, key string) error {
	_, err := rd.client.Del(ctx, key).Result()
	if (err != nil) {
		return err
	}
	return nil
}

func (rd *Cache) HasKey(ctx context.Context, key string) (int64, error) {
	res, err := rd.client.Exists(ctx, key).Result()
	if (err != nil) {
		return 0, err
	}
	return res, nil
}