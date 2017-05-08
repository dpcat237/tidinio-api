package device_model

import (
	"github.com/jinzhu/gorm"
)

const DeviceKey = "device_key"
const DeviceTable = "device"

type Device struct {
	gorm.Model

	DeviceKey string
	UserId    uint
	PushId    string
}

func (Device) TableName() string {
	return DeviceTable
}
