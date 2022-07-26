package rdb

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/go-redis/redis/v7"
	"io"
	"sync/atomic"
	"time"
)

var ErrNotObtained = errors.New("lock: not obtained")

var lockRelease = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)

// todo 後續導入redsync套件，針對globallock應用更全面
func (client *RedisClient) GlobalLock(ctx context.Context, key string, value interface{}, expire time.Duration, retry RetryStrategy) error {
	if len(key) < 0 {
		return errors.New("key can not be empty")
	}

	// make sure we don't retry forever
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(expire))
		defer cancel()
	}

	var timer *time.Timer
	for {
		isOk, err := client.Conn.SetNX(key, value, expire).Result()
		if err != nil {
			return err
		} else if isOk {
			return nil
		}

		var timeInterVal time.Duration
		if retry != nil {
			timeInterVal = retry.TimeInterval()
		}

		if timeInterVal < 1 {
			return ErrNotObtained
		}

		if timer == nil {
			timer = time.NewTimer(retry.TimeInterval())
			defer timer.Stop()
		} else {
			timer.Reset(retry.TimeInterval())
		}

		select {
		case <-ctx.Done():
			return ErrNotObtained
		case <-timer.C:
		}
	}
}

//解锁
func (client *RedisClient) GlobalUnlock(ctx context.Context, key string, value string) error {
	_, err := lockRelease.Run(client.Conn, []string{key}, value).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func (client *RedisClient) ForceUnlock(ctx context.Context, key string) error {
	if len(key) < 0 {
		return errors.New("key must not be empty")
	}

	return client.Conn.Del(key).Err()
}

func (client *RedisClient) IsOnLock(ctx context.Context, key string) bool {
	if len(key) < 0 {
		return false
	}

	result, err := client.Conn.Exists(key).Result()
	if err != nil {
		//client.log.WithFields(
		//	logrus.Fields{
		//		"command": "Exists",
		//		"key":     key,
		//	}).Error(err)
		return false
	}
	return result == 1
}

func (client *RedisClient) RandomToken() (string, error) {
	client.Lock()
	defer client.Unlock()

	temp := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, temp); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(temp), nil
}

type RetryStrategy interface {
	TimeInterval() time.Duration
}

type LimitRetry struct {
	max          int64         // 最大retry次數
	count        int64         // 紀錄retry次數
	timeInterval time.Duration // retry時間間隔
}

func NewLimitRetry(max int64, timeInterval time.Duration) RetryStrategy {
	return &LimitRetry{max: max, timeInterval: timeInterval}
}

func (r *LimitRetry) TimeInterval() time.Duration {
	if atomic.LoadInt64(&r.count) >= r.max {
		return 0
	}
	atomic.AddInt64(&r.count, 1)
	return r.timeInterval
}

type NoRetry struct {
	time.Duration
}

func (r NoRetry) TimeInterval() time.Duration {
	return 0
}

func NewNoRetry() RetryStrategy {
	return NoRetry{}
}
