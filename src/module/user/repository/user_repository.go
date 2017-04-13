package user_repository

import (
	"database/sql"
	"github.com/tidinio/src/module/user/model"
	"github.com/tidinio/src/component/repository"
	"github.com/tidinio/src/module/device/model"
)

func SaveUser(user user_model.UserBasic) {
	if (app_repository.Conn.NewRecord(user)) {
		app_repository.Conn.Create(&user)
	} else {
		app_repository.Conn.Save(&user)
	}
}

func GetUserByDeviceKey(deviceKey string) user_model.UserBasic {
	user := user_model.UserBasic{}
	app_repository.Conn.Joins("left join device on device.user_id = user.id").Where("device."+device_model.DeviceKey+" = ?", deviceKey).First(&user)

	return user
}

func GetUserByID(id int) user_model.UserBasic {
	user := user_model.UserBasic{}
	app_repository.Conn.Where("id = ?", id).First(&user)

	return user
}

func GetUsers() (*sql.Rows) {
	rows, _ := app_repository.Conn.Table("user").Rows()

	return rows
}
