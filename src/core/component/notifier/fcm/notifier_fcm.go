package app_fcm

import (
	"github.com/NaySoftware/go-fcm"
	"github.com/tidinio/src/core/device/handler"
	"github.com/tidinio/src/core/component/configuration"
)

const (
	SyncFeeds = "sync_feeds"
)

func RequireToSync(command string, userId uint) {
	devices := device_handler.GetDevicesKeyByUserId(userId)
	if (len(devices) < 1) {
		return
	}

	data := map[string]string{
		"msg": "sync_required",
		"sum": command,
	}
	serverKey, _ := app_conf.Data.String("fcm.serverKey")
	c := fcm.NewFcmClient(serverKey)
	c.NewFcmRegIdsMsg(devices, data)

	/*status, err := c.Send()
	if err == nil {
		status.PrintResults()

		return
	}

	app_logger.Error(err)*/
}
