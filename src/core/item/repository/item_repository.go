package item_repository

import (
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/item/model"
	"fmt"
)

const itemTable  = item_model.ItemTable

func GetItemByLink(repo common_repository.Repository, link string) item_model.Item {
	item := item_model.Item{}
	repo.DB.Where("link = ?", link).First(&item)

	return item
}

func GetItemsByIds(repo common_repository.Repository, ids []string) []item_model.Item {
	items := []item_model.Item{}
	repo.DB.
		Table(itemTable).
		Select("item.id, item.feed_id, item.title, item.link, item.content, item.created_at, item.published_at, feed.language").
		Joins("inner join feed on item.feed_id = feed.id and item.id IN(?)", ids).
		Scan(&items)

	return items
}

func GetLastItems(repo common_repository.Repository, feedId uint, limit int) []item_model.Item {
	items := []item_model.Item{}
	repo.DB.Table(itemTable).Where("feed_id = ?", feedId).Limit(limit).Scan(&items)

	return items
}

func GetReadItems(repo common_repository.Repository, userId uint, unreadIds []string) []item_model.Item {
	results := []item_model.Item{}
	repo.DB.
		Table(itemTable).
		Select(itemTable + ".id").
		Joins(
		"inner join " + userItemTable + " on " + userItemTable + ".item_id = " + itemTable + ".id and " + userItemTable + ".unread = ? " +
			"and " + userItemTable + ".user_id = ? AND " + userItemTable + ".id IN(?)",
		0, userId, unreadIds).
		Scan(&results)

	return results
}

func SaveItem(repo common_repository.Repository, item *item_model.Item) {
	if (repo.DB.NewRecord(item)) {
		repo.DB.Create(&item)
	} else {
		repo.DB.Save(&item)
	}
}

func SaveSharedItem(repo common_repository.Repository, item item_model.Item) item_model.Item {
	query := fmt.Sprintf(
		"INSERT INTO " + itemTable + " (title, link) VALUES('%s', '%s');",
		item.Title, item.Link)
	repo.DB.Exec(query)

	return GetItemByLink(repo, item.Link)
}
