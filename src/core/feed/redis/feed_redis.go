package feed_redis

import (
	"github.com/tidinio/src/core/component/redis"
	"github.com/tidinio/src/core/component/repository"
)

const feedKey = "feed"
const dataHash = "data_hash"

func GetDataHash(feedId uint) string {
	redis := common_redis.InitConnection()
	key := feedKey+"_"+common_repository.UintToString(feedId)+"_"+dataHash

	return redis.Get(key).String()
}

func SetDataHash(feedId uint, hash string) {
	redis := common_redis.InitConnection()
	key := feedKey+"_"+common_repository.UintToString(feedId)+"_"+dataHash
	redis.Set(key, hash, common_redis.TwoWeeks)
}
