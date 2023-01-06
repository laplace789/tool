package rdb

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetLock(lock *RedisLock) (bool, error)
	ReleaseLock(Lock *RedisLock) (bool, error)
}

func NewRedisClient(config *Config) RedisClient {
	rdb := &RedisClientImp{
		conf: config,
	}
	rdb.init()
	return rdb
}

type RedisClientImp struct {
	Conn *redis.Client
	conf *Config
	log  *logrus.Logger
	sync.Mutex
}
