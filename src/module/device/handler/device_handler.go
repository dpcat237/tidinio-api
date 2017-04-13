package device_handler

import "github.com/tidinio/src/module/device/repository"

func GetDevicesKeyByUserId(userId uint) []string {
	keys := []string{}
	devices := device_repository.GetDevicesByUserId(userId)
	for _, device := range devices {
		keys = append(keys, device.DeviceKey)
	}

	return keys
}
