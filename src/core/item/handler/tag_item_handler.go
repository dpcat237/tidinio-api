package item_handler

import (
	"github.com/tidinio/src/core/item/model"
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/item/repository"
	"github.com/tidinio/src/core/component/helper"
	"github.com/tidinio/src/core/item/data_transformer"
)

var tagItemsAdd = []item_model.TagItem{}
var tagItemsRemove = make(map[uint][]uint)

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

func SyncTagItems(apiTagItems []item_model.TagItemList) []item_model.TagItemSync {
	if (len(apiTagItems) < 1) {
		return []item_model.TagItemSync{}
	}

	updateLocalTagItems(apiTagItems)
	//TODO: executeCrawling
	userItemIds := helper_collection.GetUserItemIdsFromTagItemListCollectionStr(apiTagItems)
	repo := common_repository.InitConnection()

	return item_transformer.ToTagItemSync(item_repository.GetUnreadTagItemsSync(repo, userItemIds))
}

func addItemsTagsRelation(repo common_repository.Repository, collection []item_model.TagItem) {
	relatedTags := item_repository.GetTagsByUserItemIds(repo, helper_collection.GetUserItemIdsFromTagItemCollectionStr(collection), 0)
	if (len(relatedTags) < 1) {
		item_repository.CreateTagItems(repo, tagItemsAdd)

		return
	}

	markAsRead := make(map[uint][]uint)
	for _, relatedTag := range relatedTags {
		for itemKey, item := range collection {
			if (relatedTag.UserItemId == item.UserItemId && relatedTag.TagId == item.TagId) {
				markAsRead[item.UserItemId] = append(markAsRead[item.UserItemId], item.TagId)
				collection = append(collection[:itemKey], collection[itemKey + 1:]...)
			}
		}
	}

	if (len(relatedTags) > 0) {
		item_repository.MarkAsUnread(repo, markAsRead, 1)
	}
	item_repository.CreateTagItems(repo, tagItemsAdd)
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

func checkTagItemDifferences(tagItemList item_model.TagItemList, userItemTags []item_model.TagItem) {
	dbTags := []uint{}
	for _, userItemTag := range userItemTags {
		if (!tagItemList.HasTag(userItemTag.TagId)) {
			tagItemsRemove[userItemTag.UserItemId] = append(tagItemsRemove[userItemTag.UserItemId], userItemTag.TagId)
		}
		dbTags = append(dbTags, userItemTag.TagId)
	}

	for _, tagId := range tagItemList.Tags {
		if (!helper_collection.CheckIdExistInCollection(tagId, dbTags)) {
			tagItem := item_model.TagItem{}
			tagItem.UserItemId = tagItemList.ArticleId
			tagItem.TagId = tagId
			tagItem.Unread = 1

			tagItemsAdd = append(tagItemsAdd, tagItem)
		}
	}
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
	relatedTags := item_repository.GetTagsByUserItemIds(repo, helper_collection.GetUserItemIdsFromTagItemCollectionStr(tagItems), 1)
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

func updateLocalTagItems(tagItems []item_model.TagItemList) {
	tagItemsAdd = []item_model.TagItem{}
	tagItemsRemove = make(map[uint][]uint)
	repo := common_repository.InitConnection()

	apiItems := helper_collection.MoveTagItemsListUnderUserItemId(tagItems)
	item_repository.UpdateUserItemsStaredStatus(repo, tagItems)
	dbItems := helper_collection.
	MoveTagItemsUnderUserItemId(item_repository.GetTagsByUserItemIds(repo, helper_collection.GetUserItemIdsFromTagItemListCollectionStr(tagItems), 1))

	for userItemId, userItemTags := range dbItems {
		if (len(apiItems[userItemId].Tags) > 0) {
			checkTagItemDifferences(apiItems[userItemId], userItemTags)
		}
	}

	if (len(tagItemsRemove) > 0) {
		item_repository.MarkAsUnread(repo, tagItemsRemove, 0)
		tagItemsRemove = make(map[uint][]uint)
	}

	if (len(tagItemsAdd) > 0) {
		addItemsTagsRelation(repo, tagItemsAdd)
		tagItemsAdd = []item_model.TagItem{}
	}
}
