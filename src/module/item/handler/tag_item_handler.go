package item_handler

import (
	"github.com/tidinio/src/module/item/model"
	"github.com/tidinio/src/component/helper/collection"
	"github.com/tidinio/src/module/item/repository"
	"github.com/tidinio/src/module/item/data_transformer"
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

	tagsIdsStr := collection_helper.ConvertIntToStringSlice(tagsIds)
	unreadTagItemsIdsStr := collection_helper.ConvertIntToStringSlice(unreadTagItemsIds)

	totalUnread := item_repository.TotalUnreadTagItems(tagsIdsStr)
	tagItemsList := getUnreadTagItemsRecursive(tagsIdsStr, unreadTagItemsIdsStr, 0, limit + 5, totalUnread)
	if (len(tagItemsList) < 1) {
		return []item_model.TagItemList{}
	}

	unreadUserItemsIds := collection_helper.GetUserItemIdsFromTagItemListCollectionStr(tagItemsList)
	itemsContent := item_repository.GetUnreadUserItemsContent(unreadUserItemsIds)

	return item_model.AddTagItemsListContent(tagItemsList, itemsContent)
}

func SyncTagItems(apiTagItems []item_model.TagItemList) []item_model.TagItemSync {
	if (len(apiTagItems) < 1) {
		return []item_model.TagItemSync{}
	}

	updateLocalTagItems(apiTagItems)
	//TODO: executeCrawling
	userItemIds := collection_helper.GetUserItemIdsFromTagItemListCollectionStr(apiTagItems)

	return item_transformer.ToTagItemSync(item_repository.GetUnreadTagItemsSync(userItemIds))
}

func addItemsTagsRelation(collection []item_model.TagItem) {
	relatedTags := item_repository.GetTagsByUserItemIds(collection_helper.GetUserItemIdsFromTagItemCollectionStr(collection), 0)
	if (len(relatedTags) < 1) {
		item_repository.CreateTagItems(tagItemsAdd)

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
		item_repository.MarkAsUnread(markAsRead, 1)
	}
	item_repository.CreateTagItems(tagItemsAdd)
}

func addSharedItem(userId uint, tagId uint, title string, link string) {
	item := item_repository.GetItemByLink(link)
	if (item.ID == 0) {
		createSharedItem(userId, tagId, title, link)

		return
	}

	userItem := item_repository.GetUserItemByItemUserId(item.ID, userId)
	if (userItem.ID == 0) {
		userItem := createUserItem(userId, item.ID)
		createTagItem(userItem.ID, tagId)

		return
	}

	tagItem := item_repository.GetTagItemByUserItemTagId(userItem.ID, tagId)
	if (tagItem.ID == 0) {
		createTagItem(userItem.ID, tagId)

		return
	}

	if (!tagItem.IsUnread()) {
		tagItem.Unread = 1
		item_repository.SaveTagItem(&tagItem)
	}
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
		if (!collection_helper.CheckIdExistInCollection(tagId, dbTags)) {
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

	item = item_repository.SaveSharedItem(item)
	userItem := createUserItem(userId, item.ID)
	createTagItem(userItem.ID, tagId)
}

func createTagItem(userItemId uint, tagId uint) {
	tagItem := item_model.TagItem{}
	tagItem.TagId = tagId
	tagItem.UserItemId = userItemId
	tagItem.Unread = 1

	item_repository.SaveTagItem(&tagItem)
}

func getUnreadTagItemsRecursive(
tagsIds []string,
unreadTagItemsIds []string,
offset int,
limit int,
totalUnread int) []item_model.TagItemList {
	tagItems := item_repository.GetUnreadTagItems(tagsIds, unreadTagItemsIds, offset, limit)
	relatedTags := item_repository.GetTagsByUserItemIds(collection_helper.GetUserItemIdsFromTagItemCollectionStr(tagItems), 1)
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
	moreUnreadItems := getUnreadTagItemsRecursive(tagsIds, unreadTagItemsIds, offset, limit, totalUnread)

	return item_model.MergeTagItemsList(tagItemsList, moreUnreadItems)
}

func updateLocalTagItems(tagItems []item_model.TagItemList) {
	tagItemsAdd = []item_model.TagItem{}
	tagItemsRemove = make(map[uint][]uint)

	apiItems := collection_helper.MoveTagItemsListUnderUserItemId(tagItems)
	item_repository.UpdateUserItemsStaredStatus(tagItems)
	dbItems := collection_helper.
	MoveTagItemsUnderUserItemId(item_repository.GetTagsByUserItemIds(collection_helper.GetUserItemIdsFromTagItemListCollectionStr(tagItems), 1))

	for userItemId, userItemTags := range dbItems {
		if (len(apiItems[userItemId].Tags) > 0) {
			checkTagItemDifferences(apiItems[userItemId], userItemTags)
		}
	}

	if (len(tagItemsRemove) > 0) {
		item_repository.MarkAsUnread(tagItemsRemove, 0)
		tagItemsRemove = make(map[uint][]uint)
	}

	if (len(tagItemsAdd) > 0) {
		addItemsTagsRelation(tagItemsAdd)
		tagItemsAdd = []item_model.TagItem{}
	}
}
