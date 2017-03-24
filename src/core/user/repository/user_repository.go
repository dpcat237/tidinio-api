package user_repository

import (
	"github.com/jinzhu/gorm"
	"database/sql"
	"github.com/tidinio/src/core/user/model"
	"github.com/tidinio/src/core/device/model"
	"github.com/tidinio/src/core/component/repository"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	userRepo := UserRepository{}
	userRepo.db = common_repository.InitConnection()

	return userRepo
}

func SaveUser(userRepo UserRepository, user user_model.UserBasic) {
	if (userRepo.db.NewRecord(user)) {
		userRepo.db.Create(&user)
	} else {
		userRepo.db.Save(&user)
	}
}

func GetUserByDeviceKey(userRepo UserRepository, deviceKey string) user_model.UserBasic {
	user := user_model.UserBasic{}
	userRepo.db.Joins("left join device on device.user_id = user.id").Where("device."+device_model.DeviceKey+" = ?", deviceKey).First(&user)

	return user
}

func GetUserByID(userRepo UserRepository, id int) user_model.UserBasic {
	user := user_model.UserBasic{}
	userRepo.db.Where("id = ?", id).First(&user)

	return user
}

func GetUsers(userRepo UserRepository) (*sql.Rows) {
	rows, _ := userRepo.db.Table("user").Rows()

	return rows
}
