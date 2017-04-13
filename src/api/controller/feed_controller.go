package controller

import (
	"net/http"
	"github.com/tidinio/src/core/component/controller"
	"github.com/tidinio/src/core/feed/handler"
	"github.com/gorilla/mux"
	"github.com/tidinio/src/core/feed/model"
	"github.com/tidinio/src/core/component/helper/string"
)

type addFeed struct {
	FeedUrl string `json:"feed_url"`
}

type editFeed struct {
	FeedTitle string `json:"title"`
}

type syncFeeds struct {
	Feeds []feed_model.UserFeedSync `json:"feeds"`
}

func AddFeed(w http.ResponseWriter, r *http.Request) {
	data := addFeed{}
	user, err := app_controller.GetAuthContent(w, r, &data)
	if err != nil {
		return
	}

	userFeed, err := feed_handler.AddFeed(user.ID, data.FeedUrl)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, "Wrong url")

		return
	}
	app_controller.ReturnJson(w, userFeed)
}

func DeleteFeed(w http.ResponseWriter, r *http.Request) {
	data := editFeed{}
	user, err := app_controller.GetAuthContent(w, r, &data)
	if err != nil {
		return
	}

	vars := mux.Vars(r)
	feed_handler.UnsubscribeFromFeed(user.ID, string_helper.StringToUint(vars["id"]))
	app_controller.ReturnNoContent(w)
}

func EditFeed(w http.ResponseWriter, r *http.Request) {
	data := editFeed{}
	user, err := app_controller.GetAuthContent(w, r, &data)
	if err != nil {
		return
	}

	vars := mux.Vars(r)
	err = feed_handler.EditFeedTitle(user.ID, string_helper.StringToUint(vars["id"]), data.FeedTitle)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
	}
	app_controller.ReturnNoContent(w)
}

func SyncFeeds(w http.ResponseWriter, r *http.Request) {
	data := syncFeeds{}
	user, err := app_controller.GetAuthContent(w, r, &data)
	if err != nil {
		return
	}

	userFeeds := feed_handler.SyncUserFeeds(user.ID, data.Feeds)
	app_controller.ReturnJson(w, userFeeds)
}

//TODO: /feed/sources
