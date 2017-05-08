package device_handler

import (
	"github.com/tidinio/src/module/device/repository"
	"github.com/tidinio/src/module/device/model"
)

func AddPushNotificationId(deviceKey string, userId uint, pushId string) {
	device := device_model.Device{}
	device.DeviceKey = deviceKey
	device.UserId = userId
	device.PushId = pushId

	device_repository.SaveDevice(&device)
}

func GetDevicesKeyByUserId(userId uint) []string {
	keys := []string{}
	devices := device_repository.GetDevicesByUserId(userId)
	for _, device := range devices {
		keys = append(keys, device.DeviceKey)
	}

	return keys
}
