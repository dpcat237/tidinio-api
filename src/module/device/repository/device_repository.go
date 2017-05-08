package device_repository

import (
	"github.com/tidinio/src/module/device/model"
	"github.com/tidinio/src/component/repository"
)

const DeviceTable  = device_model.DeviceTable

func GetDevicesByUserId(userId uint) []device_model.Device {
	devices := []device_model.Device{}
	app_repository.Conn.Table(DeviceTable).Where("user_id = ?", userId).Scan(&devices)

	return devices
}

func SaveDevice(device *device_model.Device) {
	if app_repository.Conn.NewRecord(device) {
		app_repository.Conn.Create(&device)
	} else {
		app_repository.Conn.Save(&device)
	}
}
