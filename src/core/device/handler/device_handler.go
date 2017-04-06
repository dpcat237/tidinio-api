package device_handler

import (
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/device/repository"
)

func GetDevicesKeyByUserId(userId uint) []string {
	keys := []string{}
	repo := app_repository.InitConnection()
	devices := device_repository.GetDevicesByUserId(repo, userId)
	for _, device := range devices {
		keys = append(keys, device.DeviceKey)
	}

	return keys
}
