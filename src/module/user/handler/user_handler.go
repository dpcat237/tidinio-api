package user_handler

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/tidinio/src/component/helper/string"
	"github.com/tidinio/src/component/logger"
	"github.com/tidinio/src/component/notifier/email"
	"github.com/tidinio/src/module/device/handler"
	"github.com/tidinio/src/module/user/cache"
	"github.com/tidinio/src/module/user/model"
	"github.com/tidinio/src/module/user/repository"
)

const passwordLength = 25

func ChangePasswordHash(userSync user_model.UserSync, hash string) error {
	userId := user_cache.GetUserIdPasswordRecoveryHash(hash)
	if userId < 1 {
		return errors.New("Recovery hash time out")
	}
	user := user_repository.GetUserByEmail(userSync.Email)
	if userId != user.ID {
		return errors.New("Wrong user")
	}
	user.Password = prepareNewPassword(userSync.Password)
	user_repository.SaveUser(&user)
	return nil
}

func GetUserByDeviceId(deviceId string) user_model.UserBasic {
	return user_repository.GetUserByDeviceKey(deviceId)
}

func LoginPreviewUser(deviceKey string) {
	user := registerPreviewUser()
	device_handler.RegisterDevice(deviceKey, user.ID)
}

func LoginUser(deviceKey string, userSync user_model.UserBasic) error {
	if isUserLogged(deviceKey, userSync, false) {
		return nil
	}
	return errors.New("Wrong data")
}

func LoginSocialUser(deviceKey string, userSync user_model.UserBasic) error {
	if isUserLogged(deviceKey, userSync, true) {
		return nil
	}
	return errors.New("Wrong data")
}

func RegisterPreviewUser(userSync user_model.UserBasic) {
	user := user_repository.GetUserById(userSync.ID)
	user.Email = userSync.Email
	user.Password = prepareNewPassword(userSync.Password)
	user.Preview = 0
	user_repository.SaveUser(&user)
}

func RegisterUser(deviceKey string, userSync user_model.UserBasic) error {
	user, err := isUserExists(userSync.Email)
	if err != nil {
		return err
	}
	if userSync.Preview == 1 {
		user = registerPreviewUser()
	} else if user.ID != 0 {
		user = registerNewUser(user, userSync)
	} else {
		user = registerNewUser(user_model.User{}, userSync)
	}
	device_handler.RegisterDevice(deviceKey, user.ID)
	return nil
}

func RequestPasswordChange(userEmail string) error {
	user := user_repository.GetUserByEmail(userEmail)
	if user.ID < 1 {
		return errors.New("Email doesn't exist.")
	}
	hash := user_cache.SetPasswordRecoveryHash(user.ID)
	request := app_email.NewRequest([]string{userEmail}, "Tidinio: recovery your password")
	templateData := struct {
		URL string
	}{
		URL: "https://tidinio.com/user/password_change/" + hash,
	}
	request.SetFrom(app_email.EmailInfo)
	err := request.ParseEmailTemplate("src/module/user/template/password_recovery.html", templateData)
	if err != nil {
		return err
	}
	return request.SendEmail()
}

func SaveFeedback(userSync user_model.UserFeedbackSync) {
	userFeedback := user_model.UserFeedback{}
	userFeedback.Email = userSync.Email
	userFeedback.Title = userSync.Title
	userFeedback.Text = userSync.Text
	user_repository.SaveUserFeedback(&userFeedback)
}

func isUserExists(email string) (user_model.User, error) {
	user := user_repository.GetUserByEmail(email)
	if user.ID == 0 {
		return user_model.User{}, nil
	}
	if user.Registered == 1 {
		return user_model.User{}, errors.New("User exists.")
	}
	return user, nil
}

func isUserLogged(deviceKey string, userSync user_model.UserBasic, social bool) bool {
	if device_handler.IsLoggedDevice(deviceKey) {
		return true
	}
	user := user_repository.GetUserByEmail(userSync.Email)
	if user.ID < 1 {
		return false
	}
	if social {
		return true
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userSync.Password)); err != nil {
		return false
	}
	return true
}

func prepareNewPassword(password string) string {
	if password == "" {
		password = string_helper.GenerateRandomStringOfSize(passwordLength)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app_logger.Error("user_handler - prepareNewPassword: " + err.Error())
	}
	return string(hash)
}

func registerNewUser(user user_model.User, userSync user_model.UserBasic) user_model.User {
	user.Email = userSync.Email
	user.Password = prepareNewPassword(userSync.Password)
	user.Registered = 1
	user_repository.SaveUser(&user)
	return user
}

func registerPreviewUser() user_model.User {
	user := user_model.User{}
	user.Preview = 1
	user.Registered = 1
	user_repository.SaveUser(&user)
	return user
}
