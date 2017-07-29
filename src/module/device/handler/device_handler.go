package device_handler

import (
	"errors"

	"github.com/tidinio/src/module/device/cache"
	"github.com/tidinio/src/module/device/model"
	"github.com/tidinio/src/module/device/repository"
)

func GetDevicesKeyByUserId(userId uint) []string {
	keys := []string{}
	devices := device_repository.GetDevicesByUserId(userId)
	for _, device := range devices {
		keys = append(keys, device.DeviceKey)
	}

	return keys
}

func IsLoggedDevice(deviceKey string) bool {
	userId := device_cache.GetUserId(deviceKey)
	if userId < 1 {
		return false
	}
	return true
}

func RegisterDevice(deviceKey string, userId uint) {
	device := device_model.Device{}
	device.DeviceKey = deviceKey
	device.UserId = userId
	device_repository.SaveDevice(&device)
	saveSessionKey(deviceKey, userId)
}

func UpdatePushNotificationId(deviceKey string, pushId string) error {
	device := device_repository.GetDeviceByDeviceKey(deviceKey)
	if device.ID < 1 {
		return errors.New("Device doesn't exist")
	}
	device.PushId = pushId
	device_repository.SaveDevice(&device)
	go func() {
		device_repository.RemoveOldDevices(deviceKey, pushId)
	}()
	return nil
}

func saveSessionKey(deviceKey string, userId uint) {
	device_cache.SetUserId(deviceKey, userId)
}
