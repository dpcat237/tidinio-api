package app_redis

import (
	"github.com/go-redis/redis"
	"github.com/tidinio/src/component/configuration"
)

const oneMinute = 60000
const OneDay = oneMinute * 60 * 24
const SevenDays = OneDay * 7
const TwoWeeks = SevenDays * 2

func InitConnection() *redis.Client {
	host, _ := app_conf.Data.String("redis.host")
	return redis.NewClient(&redis.Options{
		Addr:     host+":6379",
		Password: "", // no password set
		DB:       0, // use default DB
	})
}
