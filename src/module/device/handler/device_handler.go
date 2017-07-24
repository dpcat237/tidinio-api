package device_handler

import (
	"github.com/tidinio/src/module/device/repository"
	"github.com/tidinio/src/module/device/model"
	"github.com/tidinio/src/module/user/handler"
	"errors"
)

func GetDevicesKeyByUserId(userId uint) []string {
	keys := []string{}
	devices := device_repository.GetDevicesByUserId(userId)
	for _, device := range devices {
		keys = append(keys, device.DeviceKey)
	}

	return keys
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
