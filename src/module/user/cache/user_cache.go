package user_cache

import (
	"github.com/cstockton/go-conv"

	"github.com/tidinio/src/component/redis"
	"github.com/tidinio/src/component/helper/string"
)

const (
	prefixPasswordRecovery = "user_password_recovery"
	passwordRecoveryLength = 20
)

func GetUserIdPasswordRecoveryHash(hash string) uint {
	redis := app_redis.InitConnection()
	key := prefixPasswordRecovery + "_" + hash
	return conv.Uint(redis.Get(key).String())
}

func SetPasswordRecoveryHash(userId uint) string {
	hash := string_helper.GenerateRandomStringOfSize(passwordRecoveryLength)
	redis := app_redis.InitConnection()
	userIdStr := conv.String(userId)
	key := prefixPasswordRecovery + "_" + hash
	redis.Set(key, userIdStr, app_redis.SevenDays)
	return hash
}
