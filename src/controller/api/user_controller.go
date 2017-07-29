package api_controller

import (
	"net/http"

	"github.com/tidinio/src/controller"
	"github.com/tidinio/src/module/user/handler"
	"github.com/tidinio/src/module/user/model"
)

func UserFeedback(w http.ResponseWriter, r *http.Request) {
	userFeedback := user_model.UserFeedbackSync{}
	err := app_controller.GetContent(r, &userFeedback)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}

	user_handler.SaveFeedback(userFeedback)
	app_controller.ReturnNoContent(w)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	deviceId, err := app_controller.GetAuthId(w, r)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}
	user := user_model.UserBasic{}
	err = app_controller.GetContent(r, &user)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}

	err = user_handler.LoginUser(deviceId, user)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
	} else {
		app_controller.ReturnNoContent(w)
	}
}

func UserPasswordRecovery(w http.ResponseWriter, r *http.Request) {
	user := user_model.UserSync{}
	err := app_controller.GetContent(r, &user)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}
	hash := app_controller.GetVariable(r, "hash")
	if hash == "" {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}

	err = user_handler.ChangePasswordHash(user, hash)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
	} else {
		app_controller.ReturnNoContent(w)
	}
}

func UserPasswordRecoveryRequest(w http.ResponseWriter, r *http.Request) {
	user := user_model.UserSync{}
	err := app_controller.GetContent(r, &user)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}

	err = user_handler.RequestPasswordChange(user.Email)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
	} else {
		app_controller.ReturnNoContent(w)
	}
}

func UserPreviewLogin(w http.ResponseWriter, r *http.Request) {
	deviceId, err := app_controller.GetAuthId(w, r)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}
	user_handler.LoginPreviewUser(deviceId)
	app_controller.ReturnNoContent(w)
}

func UserPreviewRegister(w http.ResponseWriter, r *http.Request) {
	user, err := app_controller.GetAuth(w, r)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}
	user_handler.RegisterPreviewUser(user)
	app_controller.ReturnNoContent(w)
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	deviceId, err := app_controller.GetAuthId(w, r)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}
	user := user_model.UserBasic{}
	err = app_controller.GetContent(r, &user)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}

	err = user_handler.RegisterUser(deviceId, user)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
	} else {
		app_controller.ReturnNoContent(w)
	}
}

func UserSocialLogin(w http.ResponseWriter, r *http.Request) {
	deviceId, err := app_controller.GetAuthId(w, r)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}
	user := user_model.UserBasic{}
	err = app_controller.GetContent(r, &user)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}

	err = user_handler.LoginSocialUser(deviceId, user)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
	} else {
		app_controller.ReturnNoContent(w)
	}
}
