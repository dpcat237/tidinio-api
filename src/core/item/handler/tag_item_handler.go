package item_handler

import (
	"github.com/tidinio/src/core/item/model"
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/item/repository"
)

func AddSharedItems(userId uint, collection []item_model.SharedItem) {
	for _, item := range collection {
		addSharedItem(userId, item.TagId, item.Title, item.Link)
	}
}

func addSharedItem(userId uint, tagId uint, title string, link string) {
	repo := common_repository.InitConnection()
	item := item_repository.GetItemByLink(repo, link)
	if (item.ID == 0) {
		createSharedItem(userId, tagId, title, link)

		return
	}

	userItem := item_repository.GetUserItemByItemUserId(repo, item.ID, userId)
	if (userItem.ID == 0) {
		userItem := createUserItem(userId, item.ID)
		createTagItem(userItem.ID, tagId)

		return
	}

	tagItem := item_repository.GetTagItemByUserItemTagId(repo, userItem.ID, tagId)
	if (tagItem.ID == 0) {
		createTagItem(userItem.ID, tagId)

		return
	}

	if (!tagItem.IsUnread()) {
		tagItem.Unread = 1
		item_repository.SaveTagItem(repo, &tagItem)
	}

	defer repo.Close()
}

func createSharedItem(userId uint, tagId uint, title string, link string) {
	item := item_model.Item{}
	item.Title = title
	item.Link = link

	repo := common_repository.InitConnection()
	item = item_repository.SaveSharedItem(repo, item)
	userItem := createUserItem(userId, item.ID)
	createTagItem(userItem.ID, tagId)
}

func createTagItem(userItemId uint, tagId uint) {
	tagItem := item_model.TagItem{}
	tagItem.TagId = tagId
	tagItem.UserItemId = userItemId
	tagItem.Unread = 1

	repo := common_repository.InitConnection()
	item_repository.SaveTagItem(repo, &tagItem)
}
