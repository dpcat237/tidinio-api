package controller

import (
	"net/http"
	"github.com/tidinio/src/core/component/controller"
	"github.com/tidinio/src/core/feed/handler"
)

type addFeed struct {
	FeedUrl string `json:"feed_url"`
}

func AddFeed(w http.ResponseWriter, r *http.Request) {
	data := addFeed{}
	user, err := app_controller.GetAuthContent(w, r, &data)
	if err != nil {
		return
	}

	feed, err := feed_handler.AddFeed(user.ID, data.FeedUrl)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, "Wrong url")
	}
	app_controller.ReturnJson(w, feed)
}
