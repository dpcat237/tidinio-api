package item_handler

import (
	"time"
	"github.com/mmcdole/gofeed"
	"github.com/tidinio/src/module/item/repository"
	"github.com/tidinio/src/component/helper/string"
	"github.com/tidinio/src/module/item/model"
	"github.com/tidinio/src/component/helper/http"
)

func CreateUpdateItem(itemData *gofeed.Item, feedId uint) bool {
	item := item_repository.GetItemByLink(itemData.Link)
	currentContentHash := string_helper.GetHashFromString(itemData.Content)
	if (item.ID < 1) {
		item = createItemFromFeed(feedId, currentContentHash, itemData)
		item_repository.SaveItem(&item)

		return true
	}

	if (currentContentHash != item.ContentHash) {
		item.Content = itemData.Content
		item.ContentHash = currentContentHash
		updatedDate := getUpdatedFeedTime(itemData)
		if (item.CreatedAt.Before(updatedDate)) {
			item.CreatedAt = updatedDate
		}
		item_repository.SaveItem(&item)
	}

	return false
}

func IsItemNeedsCrawling(item item_model.Item) bool {
	originalLength := string_helper.StringLength(string_helper.StripHtmlContent(item.Content))
	crawledLength := string_helper.StringLength(string_helper.StripHtmlContent(http_helper.GetContentFromUrl(item.Link)))
	difference := (crawledLength / originalLength) * 100

	return (difference > 120)
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
