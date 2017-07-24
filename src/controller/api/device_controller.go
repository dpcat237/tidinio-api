package api_controller

import (
	"net/http"
	"github.com/tidinio/src/controller"
	"github.com/tidinio/src/module/device/handler"
)

type addPush struct {
	PushId string `json:"push_notification_id"`
}

func UpdateNotificationId(w http.ResponseWriter, r *http.Request) {
	data := addPush{}
	_, err := app_controller.GetAuthContent(w, r, &data)
	if err != nil {
		return
	}

	err = device_handler.UpdatePushNotificationId(r.Header.Get("deviceId"), data.PushId)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
	} else {
		app_controller.ReturnNoContent(w)
	}
}
