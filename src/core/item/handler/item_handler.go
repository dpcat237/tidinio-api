package item_handler

import (
	"github.com/tidinio/src/core/component/repository"
	"github.com/mmcdole/gofeed"
	"github.com/tidinio/src/core/component/helper/string"
	"github.com/tidinio/src/core/item/repository"
	"github.com/tidinio/src/core/item/model"
	"time"
)

func CreateUpdateItem(repo app_repository.Repository, itemData *gofeed.Item, feedId uint) bool {
	item := item_repository.GetItemByLink(repo, itemData.Link)
	currentContentHash := string_helper.GetHashFromString(itemData.Content)
	if (item.ID < 1) {
		item = createItemFromFeed(feedId, currentContentHash, itemData)
		item_repository.SaveItem(repo, &item)

		return true
	}

	if (currentContentHash != item.ContentHash) {
		item.Content = itemData.Content
		item.ContentHash = currentContentHash
		updatedDate := getUpdatedFeedTime(itemData)
		if (item.CreatedAt.Before(updatedDate)) {
			item.CreatedAt = updatedDate
		}
		item_repository.SaveItem(repo, &item)
	}

	return false
}

func createItemFromFeed(feedId uint, contentHash string, itemData *gofeed.Item) item_model.Item {
	item := item_model.Item{}
	item.FeedId = feedId
	item.Title = itemData.Title
	item.Link = itemData.Link
	item.Content = itemData.Content
	item.ContentHash = contentHash
	item.PublishedAt = getPublishedFeedTime(itemData)
	if (itemData.Author != nil) {
		item.Author = itemData.Author.Name
	}

	return item
}

func getPublishedFeedTime(itemData *gofeed.Item) time.Time {
	if (itemData.Published != "") {
		return *itemData.PublishedParsed
	}

	return getUpdatedFeedTime(itemData)
}

func getUpdatedFeedTime(itemData *gofeed.Item) time.Time {
	if (itemData.Updated != "") {
		return *itemData.UpdatedParsed
	}

	return time.Now()
}
