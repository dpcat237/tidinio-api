package app_fcm

import (
	"github.com/NaySoftware/go-fcm"
	"github.com/tidinio/src/module/device/handler"
	"github.com/tidinio/src/component/configuration"
)

const (
	NoticeApi = "api"
	NoticeWeb = "web"

	//Feed
	SyncFeeds = "sync_feeds"

	//Filter
	AddFilters    = "add_filters"
	DeleteFilters = "delete_filters"
	UpdateFilters = "update_filters"

	//Tag
	AddTags    = "add_tags"
	DeleteTags = "delete_tags"
	UpdateTags = "update_tags"
)

func RequireToSync(command string, userId uint) {
	devices := device_handler.GetDevicesKeyByUserId(userId)
	if len(devices) < 1 {
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

func Send(userId uint, command string, data string, noticeType string) {
	devices := device_handler.GetDevicesKeyByUserId(userId)
	if !isEnoughDevices(len(devices), noticeType) {
		return
	}

	body := map[string]string{
		"msg": command,
		"sum": data,
	}
	serverKey, _ := app_conf.Data.String("fcm.serverKey")
	c := fcm.NewFcmClient(serverKey)
	c.NewFcmRegIdsMsg(devices, body)

	/*status, err := c.Send()
	if err == nil {
		status.PrintResults()

		return
	}
	app_logger.Error(err)*/
}

func isEnoughDevices(quantity int, noticeType string) bool {
	if noticeType == NoticeApi && quantity > 2 {
		return true
	} else if noticeType == NoticeWeb && quantity > 1 {
		return true
	}
	return false
}
