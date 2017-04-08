package app_fcm

import (
	"github.com/NaySoftware/go-fcm"
	"github.com/tidinio/src/core/device/handler"
)

const (
	serverKey = "AAAACFnF9nM:APA91bHWPqMubkkZVCxwh17rk-zV7L4NtNE7xwXleh41nXlzv-myuOh8eZ9i88pHObNo39CVgfmPE1sC1yciih_Np_lXNKnVXZsuQ6heEx_DmRhIvHe6-T-eSVrOsynnOHDa0awN_N5I"
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
	c := fcm.NewFcmClient(serverKey)
	c.NewFcmRegIdsMsg(devices, data)

	/*status, err := c.Send()
	if err == nil {
		status.PrintResults()

		return
	}

	app_logger.Error(err)*/
}
