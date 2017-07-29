package api_controller

import (
	"net/http"
	"github.com/tidinio/src/controller"
	"github.com/tidinio/src/module/item/model"
	"github.com/tidinio/src/module/item/handler"
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

type syncTagItems struct {
	SavedArticles []item_model.TagItemList `json:"saved_articles"`
}

func AddSharedItem(w http.ResponseWriter, r *http.Request) {
	data := sharedItems{}
	user, err := app_controller.GetAuthContent(w, r, &data)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}
	if len(data.Articles) < 1 {
		app_controller.ReturnPreconditionFailed(w, "No articles")
	}

	item_handler.AddSharedItems(user.ID, data.Articles)
	app_controller.ReturnNoContent(w)
}

func ListTagItems(w http.ResponseWriter, r *http.Request) {
	data := listTagItems{}
	_, err := app_controller.GetAuthContent(w, r, &data)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}

	items := item_handler.GetUnreadTagItems(data.SavedArticles, data.ReturnTags, data.Limit)
	app_controller.ReturnJson(w, items)
}

func SyncItems(w http.ResponseWriter, r *http.Request) {
	data := syncItems{}
	user, err := app_controller.GetAuthContent(w, r, &data)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}

	items := item_handler.SyncItems(user.ID, data.Articles, data.Limit)
	app_controller.ReturnJson(w, items)
}

func SyncTagItems(w http.ResponseWriter, r *http.Request) {
	data := syncTagItems{}
	_, err := app_controller.GetAuthContent(w, r, &data)
	if err != nil {
		app_controller.ReturnPreconditionFailed(w, err.Error())
		return
	}

	items := item_handler.SyncTagItems(data.SavedArticles)
	app_controller.ReturnJson(w, items)
}
