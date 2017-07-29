package api_controller

import (
	"net/http"

	"github.com/tidinio/src/component/notifier/fcm"
	"github.com/tidinio/src/controller"
	"github.com/tidinio/src/module/filter/handler"
	"github.com/tidinio/src/module/filter/model"
)

var filtersSync = []filter_model.FilterSync{}

func AddFilters(w http.ResponseWriter, r *http.Request) {
	user, err := app_controller.GetAuthContent(w, r, &filtersSync)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}
	filters := filter_handler.AddFilters(user.ID, filtersSync, app_fcm.NoticeApi)
	app_controller.ReturnJson(w, filters)
}

func DeleteFilters(w http.ResponseWriter, r *http.Request) {
	user, err := app_controller.GetAuthContent(w, r, &filtersSync)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}
	filter_handler.DeleteFilters(user.ID, filtersSync, app_fcm.NoticeApi)
	app_controller.ReturnNoContent(w)
}

func GetFilters(w http.ResponseWriter, r *http.Request) {
	user, err := app_controller.GetAuth(w, r)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}
	app_controller.ReturnJson(w, filter_handler.GetFilters(user.ID))
}

func UpdateFilters(w http.ResponseWriter, r *http.Request) {
	user, err := app_controller.GetAuthContent(w, r, &filtersSync)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}
	filter_handler.UpdateFilters(user.ID, filtersSync, app_fcm.NoticeApi)
	app_controller.ReturnNoContent(w)
}
