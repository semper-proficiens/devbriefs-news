package datastore

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

const expirationTTL = time.Hour * 24 // hours

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(c *redis.Client) *RedisCache {
	return &RedisCache{
		client: c,
	}
}

func (c *RedisCache) Set(key string, value any) error {
	return c.client.Set(context.TODO(), key, value, expirationTTL).Err()
}

func (c *RedisCache) Get(key string) (string, error) {

	return c.client.Get(context.TODO(), key).Result()
}

func (c *RedisCache) Remove(key string) error {
	return c.client.Del(context.TODO(), key).Err()
}

// Scan iterates over every key in the cache. Use only for debugging
func (c *RedisCache) Scan() error {
	log.Println("executing scan")
	iter := c.client.Scan(context.TODO(), 0, "*", 0).Iterator()
	for iter.Next(context.TODO()) {
		fmt.Println("keys", iter.Val())
	}
	if err := iter.Err(); err != nil {
		log.Println("error iterating keys:", err)
	}
	return nil
}

// DeleteAll all keys in the currently selected database
func (c *RedisCache) DeleteAll() error {
	return c.client.FlushDB(context.TODO()).Err()
}
