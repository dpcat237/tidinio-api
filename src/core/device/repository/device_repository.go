package device_repository

import (
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/device/model"
)

const DeviceTable  = device_model.DeviceTable

func GetDevicesByUserId(userId uint) []device_model.Device {
	devices := []device_model.Device{}
	app_repository.Conn.Table(DeviceTable).Where("user_id = ?", userId).Scan(&devices)

	return devices
}
