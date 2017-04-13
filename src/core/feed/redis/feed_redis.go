package feed_redis

import (
	"github.com/tidinio/src/core/component/redis"
	"github.com/cstockton/go-conv"
)

const feedKey = "feed"
const dataHash = "data_hash"

func GetDataHash(feedId uint) string {
	redis := app_redis.InitConnection()
	feedIdStr, _ := conv.String(feedId)
	key := feedKey+"_"+feedIdStr+"_"+dataHash

	return redis.Get(key).String()
}

func SetDataHash(feedId uint, hash string) {
	redis := app_redis.InitConnection()
	feedIdStr, _ := conv.String(feedId)
	key := feedKey+"_"+feedIdStr+"_"+dataHash
	redis.Set(key, hash, app_redis.TwoWeeks)
}
