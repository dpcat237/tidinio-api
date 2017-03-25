package user_repository

import (
	"database/sql"
	"github.com/tidinio/src/core/user/model"
	"github.com/tidinio/src/core/device/model"
	"github.com/tidinio/src/core/component/repository"
)

func SaveUser(repo common_repository.Repository, user user_model.UserBasic) {
	if (repo.DB.NewRecord(user)) {
		repo.DB.Create(&user)
	} else {
		repo.DB.Save(&user)
	}
}

func GetUserByDeviceKey(repo common_repository.Repository, deviceKey string) user_model.UserBasic {
	user := user_model.UserBasic{}
	repo.DB.Joins("left join device on device.user_id = user.id").Where("device."+device_model.DeviceKey+" = ?", deviceKey).First(&user)

	return user
}

func GetUserByID(repo common_repository.Repository, id int) user_model.UserBasic {
	user := user_model.UserBasic{}
	repo.DB.Where("id = ?", id).First(&user)

	return user
}

func GetUsers(repo common_repository.Repository) (*sql.Rows) {
	rows, _ := repo.DB.Table("user").Rows()

	return rows
}
