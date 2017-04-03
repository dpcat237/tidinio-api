package feed_redis

import (
	"github.com/tidinio/src/core/component/redis"
	"github.com/tidinio/src/core/component/repository"
)

const feedKey = "feed"
const dataHash = "data_hash"

func GetDataHash(feedId uint) string {
	redis := app_redis.InitConnection()
	key := feedKey+"_"+app_repository.UintToString(feedId)+"_"+dataHash

	return redis.Get(key).String()
}

func SetDataHash(feedId uint, hash string) {
	redis := app_redis.InitConnection()
	key := feedKey+"_"+app_repository.UintToString(feedId)+"_"+dataHash
	redis.Set(key, hash, app_redis.TwoWeeks)
}
