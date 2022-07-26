package rdb

import (
	"fmt"
	"time"
)

// redis key
const (
	GameList RedisKey = "game:list"

	VendorListKey     RedisKey = "vendor:list"      // type:hash, ex: 1:{id:1,name:xxx},2:{id:2,name:xxx}
	VendorIpKey       RedisKey = "vendor:%d:ip"     // type hash, ex: {10.10.10.1:1, 127.0.0.1:1}
	MemberKey         RedisKey = "member:%d"        //type:string ex: {id:111,vendor_id:1,nick_name:xxx...}
	GameRobotIdleKey  RedisKey = "game:robot:idle"  // type:list  ex: [123,5667,123123,....]
	GameRobotLeaveKey RedisKey = "game:robot:leave" // type:list  ex: [123,5667,123123,....]

	GlobalReportScheduleLockKey RedisKey = "global:lock:report:schedule"
	GlobalDwhScheduleLockKey    RedisKey = "global:lock:dwh:schedule"
	GlobalRobotIdleLockKey      RedisKey = "global:lock:robot:idle:%d"
	GlobalRobotFixLockKey       RedisKey = "global:lock:robot:fix"
	GlobalRobotGetLockKey       RedisKey = "global:lock:robot:get"
)

// redis ttl
var ttl = map[RedisKey]time.Duration{
	MemberKey:                   3 * 24 * time.Hour,
	GlobalReportScheduleLockKey: 30 * time.Minute,
	GlobalDwhScheduleLockKey:    30 * time.Minute,
	GlobalRobotIdleLockKey:      2 * time.Second,
	GlobalRobotFixLockKey:       10 * time.Second,
	GlobalRobotGetLockKey:       2 * time.Second,
}

type RedisKey string

func (key RedisKey) GetKey(param ...interface{}) string {
	return fmt.Sprintf(string(key), param...)
}

func (key RedisKey) GetTTL() time.Duration {
	if expired, ok := ttl[key]; ok {
		return expired
	}
	return -1
}

func (key RedisKey) ToString() string {
	return string(key)
}

func (key RedisKey) ToApmSpan(command string) string {
	return fmt.Sprintf("%s-%s", command, key)
}
