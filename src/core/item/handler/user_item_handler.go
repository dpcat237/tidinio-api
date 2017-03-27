package item_handler

import (
	"github.com/tidinio/src/core/item/model"
	"github.com/tidinio/src/core/item/repository"
	"github.com/tidinio/src/core/item/data_transformer"
	"github.com/tidinio/src/core/component/helper"
	"github.com/tidinio/src/core/component/repository"
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

	repo := common_repository.InitConnection()
	item_repository.SaveUserItem(repo, &userItem)

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
	repo := common_repository.InitConnection()
	unreadIds := helper_collection.GetIdsFromUserItemCollectionStr(collection)
	totalUnread := item_repository.CountUnreadByUser(repo, userId)
	unreadUserItems := getUnreadItemsRecursive(userId, unreadIds, 0, limit, totalUnread, repo)
	unreadItemsIds := helper_collection.GetItemIdsFromUserItemCollectionStr(unreadUserItems)
	if (len(unreadItemsIds) > 0) {
		items := item_repository.GetItemsByIds(repo, unreadItemsIds)
		userItems = mergeUserItemData(items, unreadUserItems)
	}

	if (len(unreadIds) > 0) {
		readItems := item_repository.GetReadItems(repo, userId, unreadIds)
		if (len(readItems) > 0) {
			userItems = addReadItems(userItems, unreadUserItems, readItems)
		}
	}

	return userItems
}

func getUnreadItemsRecursive(userId uint, unreadIds []string, offset int, limit int, totalUnread int, repo common_repository.Repository) []item_model.UserItem {
	unreadItems := item_repository.GetUnreadUserItems(repo, userId, unreadIds, offset, limit)
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
	moreUnreadItems := getUnreadItemsRecursive(userId, unreadIds, offset, limit, totalUnread, repo)

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
	repo := common_repository.InitConnection()
	item_repository.SyncReadItems(repo, items)
	item_repository.Close(repo)
}
