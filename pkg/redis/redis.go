package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"onboarding_test/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(r *redis.Client) *RedisClient {
	if r != nil {
		return &RedisClient{client: r}
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Envs.Host, config.Envs.RedisPort),
		Password: config.Envs.RedisPassword, // no password set
		DB:       config.Envs.RedisDB,       // use default DB
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		fmt.Printf("Error connect Redis: %v", err)
		return nil
	}
	fmt.Printf("Client: %v ", client)
	fmt.Println("Connected to Redis successfully")
	return &RedisClient{client: client}
}

func (r RedisClient) GetDriver() *redis.Client {
	return r.client
}

func (r RedisClient) GetCtx() *context.Context {
	return &ctx
}

// property is the key
// object is the value

func (r RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r RedisClient) Get(key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}