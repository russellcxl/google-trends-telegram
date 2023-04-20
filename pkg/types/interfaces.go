package types

import (
	"time"
)

type GoogleClient interface {
	GetDailyTrends(*DailyOpts) string
}

type RedisClient interface {
	SetValue(key, value string, expiry time.Duration) error
	GetValue(key string) (string, error)
	DeleteValue(key string) error 
	AddToList(key string, value string) error
	GetList(key string) ([]string, error)
	ExpireKey(key string, expiry time.Duration) error
}