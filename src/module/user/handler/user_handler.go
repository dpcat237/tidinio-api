package user_handler

import (
	"github.com/tidinio/src/module/user/model"
	"github.com/tidinio/src/module/user/repository"
)

func CreateBasicUser() {
	user := user_model.UserBasic{}
	user.Email = "email!@testy.com";
	user.Password = "testy123";

}

func GetUserByDeviceId(deviceId string) user_model.UserBasic {
	return user_repository.GetUserByDeviceKey(deviceId)
}
