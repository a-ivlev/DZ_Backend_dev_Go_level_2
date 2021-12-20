package redisDB

import (
	"context"
	"fmt"
	"rebuildCache/internal/app/cache"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
	TTL    time.Duration
}

func NewRedisClient(host, port string, ttl time.Duration) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("try to ping to redis: %w", err)
	}

	c := &RedisClient{
		Client: client,
		TTL:    ttl,
	}

	return c, nil
}

func (rc *RedisClient) Close() error {
	return rc.Client.Close()
}

var _ cache.CacheStores = &RedisClient{}

func (rc *RedisClient) Set(ctx context.Context, tagKey string, now interface{}) error {
	err := rc.Client.Set(context.Background(), tagKey, now, rc.TTL).Err()
	if err != nil {
		return fmt.Errorf("set to redis key-value: %v-%v", tagKey, now)
	}
	return nil
}

func (rc *RedisClient) Get(ctx context.Context, mkey string) ([]byte, error) {
	itemRaw, err := rc.Client.Get(context.Background(), mkey).Bytes()
	if err == redis.Nil {
		return nil, redis.Nil
	} else if err != nil {
		return nil, fmt.Errorf("redis: get info for key %v: %w", mkey, err)
	}

	return itemRaw, nil
}

func (rc *RedisClient) MGet(ctx context.Context, tags ...string) ([]interface{}, error) {
	currTags, err := rc.Client.MGet(context.Background(), tags...).Result()
	if err != nil {
		return nil, fmt.Errorf("MGET redis for tags %v: %w", tags, err)
	}
	return currTags, nil
}
