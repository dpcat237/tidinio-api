package controller

import (
	"net/http"
	"github.com/tidinio/src/core/item/model"
	"github.com/tidinio/src/core/item/handler"
	"github.com/tidinio/src/core/component/controller"
)

type listTagItems struct {
	Limit         int
	SavedArticles []uint `json:"saved_articles"`
	ReturnTags    []uint `json:"return_tags"`
}

type sharedItems struct {
	Articles []item_model.SharedItem
}

type syncItems struct {
	Limit    int
	Articles []item_model.UserItemSync
}

func AddSharedItem(w http.ResponseWriter, r *http.Request) {
	data := sharedItems{}
	user, err := common_controller.GetAuthContent(w, r, &data)
	if err != nil {
		return
	}
	if (len(data.Articles) < 1) {
		common_controller.ReturnPreconditionFailed(w, "No articles")
	}

	item_handler.AddSharedItems(user.ID, data.Articles)
	common_controller.ReturnNoContent(w)
}

func ListTagItems(w http.ResponseWriter, r *http.Request) {
	data := listTagItems{}
	_, err := common_controller.GetAuthContent(w, r, &data)
	if err != nil {
		return
	}

	items := item_handler.GetUnreadTagItems(data.SavedArticles, data.ReturnTags, data.Limit)
	common_controller.ReturnJson(w, items)
}

func SyncItems(w http.ResponseWriter, r *http.Request) {
	data := syncItems{}
	user, err := common_controller.GetAuthContent(w, r, &data)
	if err != nil {
		return
	}

	items := item_handler.SyncItems(user.ID, data.Articles, data.Limit)
	common_controller.ReturnJson(w, items)
}