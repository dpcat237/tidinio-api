package device_repository

import (
	"github.com/jinzhu/gorm"
	"github.com/tidinio/src/core/component/repository"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository() DeviceRepository {
	deviceRepo := DeviceRepository{}
	deviceRepo.db = common_repository.InitConnection()

	return deviceRepo
}

