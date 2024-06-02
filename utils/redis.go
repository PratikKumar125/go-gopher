package utils

import (
	"fmt"

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