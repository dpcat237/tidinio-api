package item_handler

import (
	"github.com/tidinio/src/core/item/model"
	"github.com/tidinio/src/core/item/repository"
	"github.com/tidinio/src/core/item/data_transformer"
	"github.com/tidinio/src/core/component/helper/collection"
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

func AddLastItemsNewUser(userId uint, feedId uint, quantity int) {
	items := item_repository.GetLastItems(feedId, quantity)
	for _, item := range items {
		userItem := item_repository.GetUserItemByItemUserId(item.ID, userId)
		if (userItem.ID < 1) {
			createUserItem(userId, item.ID)
		}
	}
}

func addReadItems(userItems []item_model.UserItemSync, unreadUserItems []item_model.UserItem, items []item_model.Item) []item_model.UserItemSync {
	for _, item := range items {
		for _, unreadUserItem := range unreadUserItems {
			if (item.ID == unreadUserItem.ItemId) {
				userItems = append(userItems, item_model.ReadToUserItemSync(unreadUserItem))
			}
		}
	}

	return userItems
}

func createUserItem(userId uint, itemId uint) item_model.UserItem {
	userItem := item_model.UserItem{}
	userItem.UserId = userId
	userItem.ItemId = itemId
	userItem.SetUnread()
	item_repository.SaveUserItem(&userItem)

	return userItem
}

func filterUnread(items []item_model.UserItem, unread bool) []item_model.UserItem {
	result := []item_model.UserItem{}
	for _, item := range items {
		if (item.IsUnread() == unread) {
			result = append(result, item)
		}
	}

	return result
}

func getUnreadItems(userId uint, collection []item_model.UserItem, limit int) []item_model.UserItemSync {
	userItems := []item_model.UserItemSync{}
	unreadIds := collection_helper.GetIdsFromUserItemCollectionStr(collection)
	totalUnread := item_repository.CountUnreadByUser(userId)
	unreadUserItems := getUnreadItemsRecursive(userId, unreadIds, 0, limit, totalUnread)
	unreadItemsIds := collection_helper.GetItemIdsFromUserItemCollectionStr(unreadUserItems)
	if (len(unreadItemsIds) > 0) {
		items := item_repository.GetItemsByIds(unreadItemsIds)
		userItems = mergeUserItemData(items, unreadUserItems)
	}

	if (len(unreadIds) > 0) {
		readItems := item_repository.GetReadItems(userId, unreadIds)
		if (len(readItems) > 0) {
			userItems = addReadItems(userItems, unreadUserItems, readItems)
		}
	}

	return userItems
}

func getUnreadItemsRecursive(userId uint, unreadIds []string, offset int, limit int, totalUnread int) []item_model.UserItem {
	unreadItems := item_repository.GetUnreadUserItems(userId, unreadIds, offset, limit)
	unreadCount := len(unreadItems)

	offset += unreadCount
	if (unreadCount >= limit || (offset + 1) >= totalUnread || limit < 5) {
		//added 5 just in case to don't do a lot of loops for few items
		return unreadItems
	}

	limit -= unreadCount
	if ((offset + limit) > totalUnread) {
		limit = totalUnread - offset
	}
	moreUnreadItems := getUnreadItemsRecursive(userId, unreadIds, offset, limit, totalUnread)

	return item_model.MergeUserItems(unreadItems, moreUnreadItems)
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

func syncReadItems(items []item_model.UserItem) {
	item_repository.SyncReadItems(items)
}
