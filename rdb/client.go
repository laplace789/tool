package rdb

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"strconv"
	"sync/atomic"
	"time"
	"tool/log"
)

func (rci *RedisClientImp) init() {
	conn := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", rci.conf.Host, rci.conf.Port),
		Username:     rci.conf.UserName,
		Password:     rci.conf.Password,
		DB:           rci.conf.Db,
		MaxRetries:   rci.conf.MaxRetries,
		PoolSize:     rci.conf.PoolSize,
		MinIdleConns: rci.conf.MinIdleConns,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			err := cn.Ping(ctx).Err()
			if err != nil {
				rci.log.WithFields(
					logrus.Fields{
						"module": "rdb",
						"type":   "connection",
						"db":     rci.conf.Db},
				).Errorf("Redis Connect Error : %v", err)
			}
			return nil
		},
	})

	rci.Conn = conn
	rci.log = logs.GetLogger()
}

func (rci *RedisClientImp) close() {
	err := rci.Conn.Close()
	if err != nil {
		rci.log.WithFields(
			logrus.Fields{
				"module": "rdb",
				"type":   "connection",
				"db":     rci.conf.Db},
		).Errorf("Redis Close Error : %v", err)
	}
}

func (rci *RedisClientImp) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := rci.Conn.Set(ctx, key, value, expiration).Err()
	if err != nil {
		rci.log.WithFields(
			logrus.Fields{
				"key":        key,
				"value":      value,
				"expiration": expiration,
			}).Errorf("Redis Command Set Error :%v", err)
		return err
	}
	return nil
}

func (rci *RedisClientImp) Get(ctx context.Context, key string) (string, error) {
	result, err := rci.Conn.Get(ctx, key).Result()

	// key 找不到一樣會噴error，所以加入redis.Nil判斷
	if err != nil && err != redis.Nil {
		rci.log.WithFields(
			logrus.Fields{
				"key": key,
			}).Errorf("Redis Command Get Error :%v", err)
		return "", err
	}

	return result, nil
}

func (rci *RedisClientImp) HMSet(ctx context.Context, key string, values map[string]interface{}) error {
	var err error
	dataMap := make(map[string]interface{})
	for k, v := range values {
		dataMap[k], err = msgpack.Marshal(v)
		if nil != err {
			return err
		}
	}

	if err := rci.Conn.HMSet(ctx, key, dataMap).Err(); err != nil {
		//Rdb.log.WithFields(
		//	logrus.Fields{
		//		"key":   key,
		//		"value": values,
		//	}).Errorf("Redis Command HMSet Error :%v", err)
		return err
	}
	return nil
}

func (rci *RedisClientImp) HMGet(ctx context.Context, key string, fields map[string]string) error {
	// 收集field
	args := make([]string, len(fields))
	i := 0
	for field := range fields {
		args[i] = field
		i++
	}

	result, err := rci.Conn.HMGet(ctx, key, args...).Result()
	if err != nil {
		rci.log.WithFields(
			logrus.Fields{
				"key":   key,
				"value": fields,
			}).Errorf("Redis Command HMSet Error :%v", err)
		return err
	}

	// map結果
	i = 0
	for _, reply := range result {
		field := args[i]
		if reply == nil {
			i++
			continue
		}
		fields[field] = reply.(string)
		i++
	}

	return nil
}

func (rci *RedisClientImp) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := rci.Conn.HGetAll(ctx, key).Result()
	if err != nil {
		rci.log.WithFields(
			logrus.Fields{
				"key":  key,
				"resp": result,
			}).Errorf("Redis Command HMSet Error :%v", err)
		return nil, err
	}
	return result, nil
}

func (rci *RedisClientImp) GetLock(lock *RedisLock) (bool, error) {
	//ctx, cancel := context.WithCancel(context.Background())
	var ctx = context.Background()
	//if lock.duration > maxLockTime {
	//	rci.log.Errorf("RedisClientImp over max lock time %d", maxLockTime)
	//	return false,errors.New("RedisClientImp over max lock time ")
	//}
	seconds := atomic.LoadUint32(&lock.duration)

	resp, err := rci.Conn.Eval(ctx, luaLockCmd, []string{lock.key}, []string{
		lock.val, strconv.Itoa(int(seconds)*millisPerSecond + tolerance),
	}).Result()

	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		rci.log.Errorf("Error on acquiring lock for %s, %s", lock.key, err.Error())
		return false, err
	} else if resp == nil {
		return false, nil
	}

	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return true, nil
	}

	rci.log.Errorf("Unknown reply when acquiring lock for %s: %v", lock.key, resp)
	return false, nil
}

func (rci *RedisClientImp) ReleaseLock(Lock *RedisLock) (bool, error) {
	var ctx = context.Background()
	resp, err := rci.Conn.Eval(ctx, luaReleaseCmd, []string{Lock.key}, []string{Lock.val}).Result()
	if err != nil {
		return false, err
	}

	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}

	return reply == 1, nil
}
