package session

import (
	"fmt"
	"log"
	"time"

	"github.com/russellcxl/google-trends/config"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

// Tries to connect to redis and returns a client if successful
func New(addr, password string, dbIndex int) (*Redis, error) {
	cfg := config.GetConfig()
	retryCount := cfg.Redis.RetryCount
	var client *redis.Client
	var err error
	for i := 0; i < retryCount; i++ {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       dbIndex,
		})
		_, err = client.Ping().Result()
		if err == nil {
			log.Println("Connected to Redis!")
			return &Redis{
				client: client,
			}, nil
		}
		log.Printf("Error connecting to Redis (%v): %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to redis after %d retries: %v", retryCount, err)
}

func (c *Redis) SetValue(key, value string, expiry time.Duration) error {
	return c.client.Set(key, value, expiry).Err()
}

func (c *Redis) GetValue(key string) (string, error) {
	return c.client.Get(key).Result()
}

func (c *Redis) KeyExists(key string) (int64, error) {
	return c.client.Exists(key).Result()
}

func (c *Redis) DeleteValue(key string) error {
	return c.client.Del(key).Err()
}

func (c *Redis) AddToList(key string, value string) error {
	return c.client.RPush(key, value).Err()
}

func (c *Redis) GetList(key string) ([]string, error) {
	return c.client.LRange(key, 0, -1).Result()
}

func (c *Redis) ExpireKey(key string, expiry time.Duration) error {
	return c.client.Expire(key, expiry).Err()
}
