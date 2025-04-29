package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Storage interface {
	Increment(key string, window time.Duration) (int64, error)
	SetExpiration(key string, duration time.Duration) error
	SetValue(key string, value int, duration time.Duration) error
	GetValue(key string) (int64, error)
}

type RedisStorage struct {
	client  *redis.Client
	context context.Context
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client:  client,
		context: context.Background(),
	}
}

func (r *RedisStorage) Increment(key string, window time.Duration) (int64, error) {

	pipeline := r.client.TxPipeline()
	incr := pipeline.Incr(r.context, key)
	pipeline.Expire(r.context, key, window)
	_, err := pipeline.Exec(r.context)
	if err != nil {
		fmt.Println("Error executing pipeline: ", err)
		return 0, err
	}

	return incr.Val(), nil
}

func (r *RedisStorage) SetExpiration(key string, duration time.Duration) error {
	return r.client.Expire(r.context, key, duration).Err()
}

func (r *RedisStorage) SetValue(key string, value int, duration time.Duration) error {
	return r.client.Set(r.context, key, value, duration).Err()
}

func (r *RedisStorage) GetValue(key string) (int64, error) {
	val, err := r.client.Get(r.context, key).Result()
	if err != nil {
		return 0, err
	}

	if val == "BLOCKED" {
		return 0, errors.New("key is blocked")
	}

	return r.client.Get(r.context, key).Int64()
}
