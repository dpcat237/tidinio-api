package api_controller

import (
	"net/http"
	"github.com/tidinio/src/controller"
	"github.com/tidinio/src/module/device/handler"
)

type addPush struct {
	PushId string `json:"push_notification_id"`
}

func AddDevice(w http.ResponseWriter, r *http.Request) {
	data := addPush{}
	user, err := app_controller.GetAuthContent(w, r, &data)
	if err != nil {
		return
	}

	device_handler.AddPushNotificationId(r.Header.Get("deviceId"), user.ID, data.PushId)
	app_controller.ReturnNoContent(w)
}
