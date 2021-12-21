package redisDB

import (
	"context"
	"fmt"
	"session-srv/internal/app/sessions"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	*redis.Client
	TTL time.Duration
}

func NewRedisClient(host, port string, ttl time.Duration) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("try to ping to redisDB: %w", err)
	}

	c := &RedisClient{
		Client: client,
		TTL:    ttl,
	}

	return c, nil
}

func (c *RedisClient) Close() error {
	return c.Client.Close()
}

var _ sessions.SessionStores = &RedisClient{}

func (c *RedisClient) SetCache(key string, value []byte) error {
	err := c.Set(context.Background(), key, value, c.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: set key %q: %w", key, err)
	}
	return nil
}

func (c *RedisClient) GetRecordCache(key string) ([]byte, error) {
	data, err := c.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	} else if data == nil {
		// add here custom err handling
		return nil, nil
	}

	return data, nil
}

func (c *RedisClient) DelCache(key string) error {
	err := c.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}

	return nil
}
