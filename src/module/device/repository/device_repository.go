package device_repository

import (
	"github.com/tidinio/src/component/repository"
	"github.com/tidinio/src/module/device/model"
)

const DeviceTable  = device_model.DeviceTable

func GetDeviceByDeviceKey(deviceKey string) device_model.Device {
	device := device_model.Device{}
	app_repository.Conn.Table(DeviceTable).Where("device_key = ?", deviceKey).First(&device)
	return device
}

func GetDevicesByUserId(userId uint) []device_model.Device {
	devices := []device_model.Device{}
	app_repository.Conn.Table(DeviceTable).Where("user_id = ?", userId).Scan(&devices)

	return devices
}

func RemoveOldDevices(deviceKey string, pushId string) {
	app_repository.Conn.Table(DeviceTable).Where("device_key = ? AND push_id != ?", deviceKey, pushId).Delete(device_model.Device{})
}

func SaveDevice(device *device_model.Device) {
	if app_repository.Conn.NewRecord(device) {
		app_repository.Conn.Create(&device)
	} else {
		app_repository.Conn.Save(&device)
	}
}
