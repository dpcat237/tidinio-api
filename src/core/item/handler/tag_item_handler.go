package item_handler

import (
	"github.com/tidinio/src/core/item/model"
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/item/repository"
	"github.com/tidinio/src/core/component/helper"
)

func AddSharedItems(userId uint, collection []item_model.SharedItem) {
	for _, item := range collection {
		addSharedItem(userId, item.TagId, item.Title, item.Link)
	}
}

func GetUnreadTagItems(unreadTagItemsIds []uint, tagsIds []uint, limit int) []item_model.TagItemList {
	if (limit < 1) {
		return []item_model.TagItemList{}
	}

	tagsIdsStr := helper_collection.ConvertIntToStringSlice(tagsIds)
	unreadTagItemsIdsStr := helper_collection.ConvertIntToStringSlice(unreadTagItemsIds)
	repo := common_repository.InitConnection()

	totalUnread := item_repository.TotalUnreadTagItems(repo, tagsIdsStr)
	tagItemsList := getUnreadTagItemsRecursive(repo, tagsIdsStr, unreadTagItemsIdsStr, 0, limit + 5, totalUnread)
	if (len(tagItemsList) < 1) {
		return []item_model.TagItemList{}
	}

	unreadUserItemsIds := helper_collection.GetUserItemIdsFromTagItemListCollectionStr(tagItemsList)
	itemsContent := item_repository.GetUnreadUserItemsContent(repo, unreadUserItemsIds)

	return item_model.AddTagItemsListContent(tagItemsList, itemsContent)
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

func getUnreadTagItemsRecursive(
repo common_repository.Repository,
tagsIds []string,
unreadTagItemsIds []string,
offset int,
limit int,
totalUnread int) []item_model.TagItemList {
	tagItems := item_repository.GetUnreadTagItems(repo, tagsIds, unreadTagItemsIds, offset, limit)
	relatedTags := item_repository.GetTagsByUserItemIds(repo, helper_collection.GetUserItemIdsFromTagItemCollectionStr(tagItems))
	tagItemsList := item_model.MergeToTagItemsList(tagItems, item_model.JoinTagsByUserItem(relatedTags))

	unreadCount := len(tagItemsList)
	offset += unreadCount
	if (unreadCount >= limit || (offset + 1) >= totalUnread || limit < 5) {
		//added 5 just in case to don't do a lot of loops for few items
		return tagItemsList
	}

	limit -= unreadCount
	if ((offset + limit) > totalUnread) {
		limit = totalUnread - offset
	}
	moreUnreadItems := getUnreadTagItemsRecursive(repo, tagsIds, unreadTagItemsIds, offset, limit, totalUnread)

	return item_model.MergeTagItemsList(tagItemsList, moreUnreadItems)
}
