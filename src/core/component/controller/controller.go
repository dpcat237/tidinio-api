package common_controller

import (
	"io/ioutil"
	"io"
	"net/http"
	"github.com/tidinio/src/core/component/model"
	"github.com/tidinio/src/core/user/handler"
	"github.com/tidinio/src/core/user/model"
	"encoding/json"
)

func GetAuthContent(w http.ResponseWriter, r *http.Request, data interface{}) (user_model.UserBasic, error)   {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, r.ContentLength))
	deviceId := r.Header.Get("deviceId")
	if err != nil || deviceId == "" {
		ReturnPreconditionFailed(w, "Authentification failed")

		return user_model.UserBasic{}, common_error.NewError("Authentification failed")
	}
	user := user_handler.GetUserByDeviceId(deviceId)
	json.Unmarshal(body, data)

	return user, err
}

func ReturnJson(w http.ResponseWriter, v interface{})  {
	json.NewEncoder(w).Encode(v)
}

func ReturnPreconditionFailed(w http.ResponseWriter, s string)  {
	http.Error(w, s, http.StatusPreconditionFailed)
}

func ReturnNoContent(w http.ResponseWriter)  {
	w.WriteHeader(http.StatusNoContent)
}
