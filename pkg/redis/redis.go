package redis

import (
	redis "github.com/go-redis/redis/v8"
)

func init() {}

func GetRedis(host string, pass string, db int) (rdb *redis.Client) {
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
		Password: pass,
		DB: db,
	})

	return rdb
}
