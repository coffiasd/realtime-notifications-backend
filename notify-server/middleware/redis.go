package middleware

import (
	"notify-server/config"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func GetRedis() *redis.Client {
	return rdb
}

func InitRedis() {
	//connect redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.ParseConfig.Redis.Addr,
		Password: config.ParseConfig.Redis.Password,
		DB:       config.ParseConfig.Redis.DB,
	})
}
