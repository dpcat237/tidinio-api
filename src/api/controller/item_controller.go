package controller

import (
	"net/http"
	"io/ioutil"
	"io"
	"encoding/json"
	"github.com/tidinio/src/core/item/model"
	"github.com/tidinio/src/core/user/handler"
	"github.com/tidinio/src/core/item/handler"
)

type syncItems struct {
	Limit int
	Articles []item_model.UserItemSync
}

func SyncItems(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, r.ContentLength))
	deviceId := r.Header.Get("deviceId")
	if err != nil || deviceId == "" {
		http.Error(w, "", http.StatusPreconditionFailed)
		return
	}

	user := user_handler.GetUserByDeviceId(deviceId)
	data := syncItems{}
	json.Unmarshal(body, &data)

	items := item_handler.SyncItems(user.ID, data.Articles, data.Limit)
	json.NewEncoder(w).Encode(items)
}
