package device_repository

import (
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/device/model"
)

func GetDevicesByUserId(repo app_repository.Repository, userId uint) []device_model.Device {
	devices := []device_model.Device{}
	repo.DB.Where("user_id = ?", userId).Scan(&devices)

	return devices
}
