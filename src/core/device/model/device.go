package device_model

import (
	"github.com/jinzhu/gorm"
	"github.com/tidinio/src/core/user/model"
)

const DeviceKey  = "device_key"
const DeviceTable = "device"

type Device struct {
	gorm.Model

	DeviceKey string `gorm:"size:255"`
	User      user_model.UserBasic
	GcmId     string `gorm:"size:255"`
}

func (Device) TableName() string {
	return DeviceTable
}
