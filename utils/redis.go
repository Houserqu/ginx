package utils

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var Redis *redis.Client

func InitRedis() {
	if viper.GetString("redis.addr") != "" {
		Redis = redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis.addr"),
			Username: viper.GetString("redis.username"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		})
	}
}
