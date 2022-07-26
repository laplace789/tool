package rdb

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"sync"
)

type RedisClient struct {
	Conn *redis.Client
	conf *Config
	log  *logrus.Logger
	sync.Mutex
}

func NewRedisClient(config *Config) *RedisClient {
	rdb := new(RedisClient)
	rdb.conf = config
	return rdb
}

func (client *RedisClient) connect() {
	conn := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", client.conf.Host, client.conf.Port),
		Username:     client.conf.UserName,
		Password:     client.conf.Password,
		DB:           client.conf.Db,
		MaxRetries:   client.conf.MaxRetries,
		PoolSize:     client.conf.PoolSize,
		MinIdleConns: client.conf.MinIdleConns,
		OnConnect: func(cn *redis.Conn) error {
			err := cn.Ping().Err()
			if err != nil {
				client.log.WithFields(
					logrus.Fields{
						"module": "redis",
						"type":   "connection",
						"db":     client.conf.Db},
				).Errorf("Redis Connect Error : %v", err)
			}

			return nil
		},
	})

	client.Conn = conn
}

func (client *RedisClient) close() {
	err := client.Conn.Close()
	if err != nil {
		client.log.WithFields(
			logrus.Fields{
				"module": "redis",
				"type":   "connection",
				"db":     client.conf.Db},
		).Errorf("Redis Close Error : %v", err)
	}
}
