package rdb

import (
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

func (client *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := client.Conn.Set(key, value, expiration).Err()
	if err != nil {
		//Rdb.log.WithFields(
		//	logrus.Fields{
		//		"key":        key,
		//		"value":      value,
		//		"expiration": expiration,
		//	}).Errorf("Redis Command Set Error :%v", err)
		return err
	}
	return nil
}

func (client *RedisClient) Get(ctx context.Context, key string) (string, error) {
	result, err := client.Conn.Get(key).Result()

	// key 找不到一樣會噴error，所以加入redis.Nil判斷
	if err != nil && err != redis.Nil {
		//Rdb.log.WithFields(
		//	logrus.Fields{
		//		"key": key,
		//	}).Errorf("Redis Command Get Error :%v", err)
		return "", err
	}

	return result, nil
}

func (client *RedisClient) HMSet(ctx context.Context, key string, values map[string]interface{}) error {
	var err error
	dataMap := make(map[string]interface{})
	for k, v := range values {
		dataMap[k], err = msgpack.Marshal(v)
		if nil != err {
			return err
		}
	}

	if err := client.Conn.HMSet(key, dataMap).Err(); err != nil {
		//Rdb.log.WithFields(
		//	logrus.Fields{
		//		"key":   key,
		//		"value": values,
		//	}).Errorf("Redis Command HMSet Error :%v", err)
		return err
	}
	return nil
}

func (client *RedisClient) HMGet(ctx context.Context, key string, fields map[string]string) error {
	// 收集field
	args := make([]string, len(fields))
	i := 0
	for field := range fields {
		args[i] = field
		i++
	}

	result, err := client.Conn.HMGet(key, args...).Result()
	if err != nil {
		//Rdb.log.WithFields(
		//	logrus.Fields{
		//		"key":   key,
		//		"value": fields,
		//	}).Errorf("Redis Command HMSet Error :%v", err)
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

func (client *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := client.Conn.HGetAll(key).Result()
	if err != nil {
		//Rdb.log.WithFields(
		//	logrus.Fields{
		//		"key":  key,
		//		"resp": result,
		//	}).Errorf("Redis Command HMSet Error :%v", err)
		return nil, err
	}
	return result, nil
}
