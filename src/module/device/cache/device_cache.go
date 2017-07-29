package device_cache

import (
	"github.com/cstockton/go-conv"

	"github.com/tidinio/src/component/redis"
)

const (
	deviceKeyHash = "device"
	userIdHash    = "user_id"
)

func GetUserId(deviceKey string) uint {
	redis := app_redis.InitConnection()
	key := deviceKeyHash + "_" + deviceKey + "_" + userIdHash
	return conv.Uint(redis.Get(key).String())
}

func SetUserId(deviceKey string, userId uint) {
	redis := app_redis.InitConnection()
	userIdStr := conv.String(userId)
	key := deviceKeyHash + "_" + deviceKey + "_" + userIdHash
	redis.Set(key, userIdStr, app_redis.TwoWeeks)
}
