package user_repository

import (
	"github.com/tidinio/src/component/repository"
	"github.com/tidinio/src/module/device/model"
	"github.com/tidinio/src/module/user/model"
)

func GetUserByDeviceKey(deviceKey string) user_model.UserBasic {
	user := user_model.UserBasic{}
	app_repository.Conn.Joins("left join device on device.user_id = user.id").
		Where("device."+device_model.DeviceKey+" = ?", deviceKey).First(&user)
	return user
}

func GetUserByEmail(email string) user_model.User {
	user := user_model.User{}
	app_repository.Conn.Where("email = ?", email).First(&user)
	return user
}

func GetUserById(id uint) user_model.User {
	user := user_model.User{}
	app_repository.Conn.Where("id = ?", id).First(&user)
	return user
}

func SaveUser(user *user_model.User) {
	if app_repository.Conn.NewRecord(user) {
		app_repository.Conn.Create(&user)
	} else {
		app_repository.Conn.Save(&user)
	}
}

func SaveUserFeedback(user *user_model.UserFeedback) {
	if app_repository.Conn.NewRecord(user) {
		app_repository.Conn.Create(&user)
	} else {
		app_repository.Conn.Save(&user)
	}
}
