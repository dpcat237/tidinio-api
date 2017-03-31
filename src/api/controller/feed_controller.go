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
	user, err := common_controller.GetAuthContent(w, r, &data)
	if err != nil {
		return
	}

	feed_handler.AddFeed(user.ID, data.FeedUrl)
}
