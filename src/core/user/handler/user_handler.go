package user_handler

import (
	"github.com/tidinio/src/core/user/model"
	"github.com/tidinio/src/core/user/repository"
	"github.com/tidinio/src/core/component/repository"
)

func CreateBasicUser() {
	user := user_model.UserBasic{}
	user.Email = "email!@testy.com";
	user.Password = "testy123";

}

func GetUserByDeviceId(deviceId string) user_model.UserBasic {
	return user_repository.GetUserByDeviceKey(app_repository.InitConnection(), deviceId)
}
