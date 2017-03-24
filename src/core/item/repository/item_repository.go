package item_repository

import (
	"github.com/jinzhu/gorm"
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/item/model"
)

const itemTable  = item_model.ItemTable

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository() ItemRepository {
	repo := ItemRepository{}
	repo.db = common_repository.InitConnection()

	return repo
}

func GetItemByIds(repo ItemRepository, ids []string) []item_model.Item {
	items := []item_model.Item{}
	repo.db.
		Table(itemTable).
		Select("item.id, item.feed_id, item.title, item.link, item.content, item.created_at, feed.language").
		Joins("inner join feed on item.feed_id = feed.id and item.id IN(?)", ids).
		Scan(&items)

	return items
}

func GetReadItems(repo ItemRepository, userId uint, unreadIds []string) []item_model.Item {
	results := []item_model.Item{}
	repo.db.
		Table(itemTable).
		Select(itemTable + ".id").
		Joins(
		"inner join " + userItemTable + " on " + userItemTable + ".item_id = " + itemTable + ".id and " + userItemTable + ".unread = ? " +
			"and " + userItemTable + ".user_id = ? AND " + userItemTable + ".id IN(?)",
		0, userId, unreadIds).
		Scan(&results)

	return results
}
