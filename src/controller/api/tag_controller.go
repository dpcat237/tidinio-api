package api_controller

import (
	"net/http"

	"github.com/tidinio/src/component/notifier/fcm"
	"github.com/tidinio/src/controller"
	"github.com/tidinio/src/module/tag/model"
	"github.com/tidinio/src/module/tag/handler"
)

var tagsSync = []tag_model.TagSync{}

func AddTags(w http.ResponseWriter, r *http.Request) {
	user, err := app_controller.GetAuthContent(w, r, &tagsSync)
	if err != nil {
		return
	}
	tags := tag_handler.AddTags(user.ID, tagsSync, app_fcm.NoticeApi)
	app_controller.ReturnJson(w, tags)
}

func DeleteTags(w http.ResponseWriter, r *http.Request) {
	user, err := app_controller.GetAuthContent(w, r, &tagsSync)
	if err != nil {
		return
	}
	tag_handler.DeleteTags(user.ID, tagsSync, app_fcm.NoticeApi)
	app_controller.ReturnNoContent(w)
}

func GetTags(w http.ResponseWriter, r *http.Request) {
	user, err := app_controller.GetAuth(w, r)
	if err != nil {
		return
	}
	app_controller.ReturnJson(w, tag_handler.GetTags(user.ID))
}

func UpdateTags(w http.ResponseWriter, r *http.Request) {
	user, err := app_controller.GetAuthContent(w, r, &tagsSync)
	if err != nil {
		return
	}
	tag_handler.UpdateTags(user.ID, tagsSync, app_fcm.NoticeApi)
	app_controller.ReturnNoContent(w)
}
