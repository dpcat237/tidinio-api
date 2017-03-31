package item_redis

import "github.com/tidinio/src/core/component/redis"

const itemKey = "item"

func GetDataHash(feedId uint) string {
	redis := common_redis.InitConnection()
	key := feedKey+"_"+feedId+"_"+dataHash

	return redis.Get(key)
}

func SetDataHash(feedId uint, hash string) {
	redis := common_redis.InitConnection()
	key := feedKey+"_"+feedId+"_"+dataHash

	return redis.Set(key, hash, common_redis.TwoWeeks)
}
