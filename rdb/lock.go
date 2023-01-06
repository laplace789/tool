package rdb

import "math/rand"

type RedisLock struct {
	duration uint32
	key      string
	val      string
}

func NewRedisLock(key string, duration uint32) *RedisLock {
	return &RedisLock{
		duration: duration,
		key:      key,
		val:      randomStr(randomLen),
	}
}

func randomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
