package session

import (
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func New(addr, password string, dbIndex int) *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       dbIndex,
		}),
	}
}

func (c *Redis) SetValue(key, value string, expiry time.Duration) error {
	return c.client.Set(key, value, expiry).Err()
}

func (c *Redis) GetValue(key string) (string, error) {
	return c.client.Get(key).Result()
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