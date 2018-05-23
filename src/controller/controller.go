package app_controller

import (
	"errors"
	"net/http"
	"io/ioutil"
	"io"
	"encoding/json"

	"github.com/gorilla/mux"

	"github.com/tidinio/src/module/user/model"
	"github.com/tidinio/src/module/user/handler"
)

func GetAuth(w http.ResponseWriter, r *http.Request) (user_model.UserBasic, error) {
	deviceKey := r.Header.Get("deviceId")
	if deviceKey == "" {
		ReturnPreconditionFailed(w, "Authentification failed")
		return user_model.UserBasic{}, errors.New("Authentification failed")
	}
	user := user_handler.GetUserByDeviceId(deviceKey)
	if user.ID < 1 {
		return user_model.UserBasic{}, errors.New("Authentification failed")
	}
	return user, nil
}

func GetAuthContent(w http.ResponseWriter, r *http.Request, data interface{}) (user_model.UserBasic, error) {
	user, err := GetAuth(w, r)
	if err != nil {
		return user, err
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, r.ContentLength))
	json.Unmarshal(body, data)
	return user, err
}

func GetAuthId(w http.ResponseWriter, r *http.Request) (string, error) {
	deviceKey := r.Header.Get("deviceId")
	if deviceKey == "" {
		return "", errors.New("No device ID.")
	}
	return deviceKey, nil
}

func GetContent(r *http.Request, data interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, r.ContentLength))
	json.Unmarshal(body, data)
	return err
}

func GetVariable(r *http.Request, key string) string {
	vars := mux.Vars(r)
	return vars[key]
}

func ReturnJson(w http.ResponseWriter, v interface{}) {
	json.NewEncoder(w).Encode(v)
}

func ReturnPreconditionFailed(w http.ResponseWriter, s string) {
	http.Error(w, s, http.StatusPreconditionFailed)
}

func ReturnNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func ReturnStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func ReturnServerFailed(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
