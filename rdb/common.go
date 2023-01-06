package rdb

const (
	luaLockCmd = `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
			return "OK"
		else
    		return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
		end
	`
	luaReleaseCmd = `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
		`
	tolerance       = 500 // milliseconds
	millisPerSecond = 1000
	maxLockTime     = 10000 // milliseconds
)

const (
	letters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomLen = 16
)
