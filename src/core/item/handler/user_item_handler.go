package item_handler

import (
	"github.com/tidinio/src/core/item/model"
	"github.com/tidinio/src/core/item/repository"
	"github.com/tidinio/src/core/item/data_transformer"
	"github.com/tidinio/src/core/component/helper"
)

func SyncItems(userId uint, collection []item_model.UserItemSync, limit int) []item_model.UserItemSync {
	if (limit < 1) {
		return []item_model.UserItemSync{}
	}

	items := item_transformer.FromUserItemsSync(collection)
	readItems := filterUnread(items, false)
	if (len(readItems) > 0) {
		syncReadItems(readItems)
	}
	unreadItems := filterUnread(items, true)

	return getUnreadItems(userId, unreadItems, limit)
}

func addReadItems (userItems []item_model.UserItemSync, unreadUserItems []item_model.UserItem, items []item_model.Item) []item_model.UserItemSync {
	for _, item := range items {
		for _, unreadUserItem := range unreadUserItems {
			if (item.ID == unreadUserItem.ItemId) {
				userItems = append(userItems, item_model.ReadToUserItemSync(unreadUserItem))
			}
		}
	}

	return userItems
}

func filterUnread(items []item_model.UserItem, unread bool) []item_model.UserItem {
	result := []item_model.UserItem{}
	for _, item := range items {
		if (item_model.IsUnread(item) == unread) {
			result = append(result, item)
		}
	}

	return result
}

func getUnreadItems(userId uint, collection []item_model.UserItem, limit int) []item_model.UserItemSync {
	userItems := []item_model.UserItemSync{}
	userItemRepo := item_repository.NewUserItemRepository()
	itemRepo := item_repository.NewItemRepository()
	unreadIds := helper_collection.GetIdsFromUserItemCollectionStr(collection)
	totalUnread := item_repository.CountUnreadByUser(userItemRepo, userId)
	unreadUserItems := getUnreadItemsRecursive(userId, unreadIds, 0, limit, totalUnread, userItemRepo)
	unreadItemsIds := helper_collection.GetItemIdsFromUserItemCollectionStr(unreadUserItems)
	if (len(unreadItemsIds) > 0) {
		items := item_repository.GetItemByIds(itemRepo, unreadItemsIds)
		userItems = mergeUserItemData(items, unreadUserItems)
	}

	if (len(unreadIds) > 0) {
		readItems := item_repository.GetReadItems(itemRepo, userId, unreadIds)
		if (len(readItems) > 0) {
			userItems = addReadItems(userItems, unreadUserItems, readItems)
		}
	}

	return userItems
}

func getUnreadItemsRecursive(userId uint, unreadIds []string, offset int, limit int, totalUnread int, userItemRepo item_repository.UserItemRepository) []item_model.UserItem {
	unreadItems := item_repository.GetUnreadUserItems(userItemRepo, userId, unreadIds, offset, limit)
	unreadCount := len(unreadItems)
	if (unreadCount < 1) {
		return unreadItems
	}

	offset += unreadCount
	if (unreadCount >= limit || (offset +1) >= totalUnread || limit < 5) { //added 5 just in case to don't do a lot of loops for few items
		return unreadItems
	}

	limit -= unreadCount
	if ((offset + limit) > totalUnread) {
		limit = totalUnread - offset
	}
	moreUnreadItems := getUnreadItemsRecursive(userId, unreadIds, offset, limit, totalUnread, userItemRepo)

	return mergeUserItems(unreadItems, moreUnreadItems)
}

func mergeUserItemData(items []item_model.Item, userItems []item_model.UserItem) []item_model.UserItemSync {
	results := []item_model.UserItemSync{}
	for _, item := range items {
		for _, userItem := range userItems {
			if (item.ID == userItem.ItemId) {
				results = append(results, item_model.MergeToUserItemSync(item, userItem))
			}
		}
	}

	return results
}

func mergeUserItems(collection1 []item_model.UserItem, collection2 []item_model.UserItem) []item_model.UserItem {
	for _, value := range collection2 {
		collection1 = append(collection1, value)
	}

	return collection1
}

func syncReadItems(items []item_model.UserItem) {
	userItemRepo := item_repository.NewUserItemRepository()
	item_repository.SyncReadItems(userItemRepo, items)
	item_repository.Close(userItemRepo)
}
