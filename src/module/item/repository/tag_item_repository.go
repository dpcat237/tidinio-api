package item_repository

import (
	"time"
	"fmt"
	"github.com/tidinio/src/module/item/model"
	"github.com/tidinio/src/component/repository"
)

const tagItemTable  = item_model.TagItemTable

func CreateTagItems(collection []item_model.TagItem) {
	tx := app_repository.Conn.Begin()
	for _, tagItem := range collection {
		tx.Create(&tagItem)
	}
	tx.Commit()
}

func GetTagItemByUserItemTagId(userItemId uint, tagId uint) item_model.TagItem {
	tagItem := item_model.TagItem{}
	app_repository.Conn.Where("user_item_id = ? AND tag_id = ?", userItemId, tagId).First(&tagItem)

	return tagItem
}

func GetTagsByUserItemIds(userItemIds []string, unread int) []item_model.TagItem {
	results := []item_model.TagItem{}
	app_repository.Conn.
		Table(tagItemTable).
		Select("id, user_item_id, tag_id").
		Where("unread = ? AND user_item_id IN(?)", unread, userItemIds).
		Scan(&results)

	return results
}

func GetUnreadTagItems(
tagsIds []string,
tagItemsIds []string,
offset int,
limit int) []item_model.TagItem {
	results := []item_model.TagItem{}
	app_repository.Conn.
		Table(tagItemTable).
		Select("id, user_item_id, unread").
		Where("unread = ? AND tag_id IN(?) AND id NOT IN(?)", 1, tagsIds, tagItemsIds).
		Offset(offset).
		Limit(limit).
		Scan(&results)

	return results
}

func GetUnreadTagItemsSync(userItemIds []string) []item_model.TagItemSyncDB {
	results := []item_model.TagItemSyncDB{}
	app_repository.Conn.
		Table(tagItemTable).
		Select(tagItemTable + ".id, " + tagItemTable + ".user_item_id, " + tagItemTable + ".tag_id, " + userItemTable +
		".stared").
		Joins(
		"inner join " + userItemTable + " on " + tagItemTable + ".user_item_id = " + userItemTable + ".id " +
			"AND " + userItemTable + ".id IN(?) AND " + tagItemTable + ".unread = ?", userItemIds, 1).
		Scan(&results)

	return results
}

func MarkAsUnread(collection map[uint][]uint, unread int) {
	tx := app_repository.Conn.Begin()
	now := time.Now().Format("2006-01-02 15:04:05")

	for userItemId, tagsId := range collection {
		for _, tagId := range tagsId {
			query := fmt.Sprintf(
				"UPDATE later_item SET unread='%d', updated_at='%s' WHERE user_item_id='%d' AND tag_id='%d';",
				unread, now, userItemId, tagId)
			tx.Exec(query)
		}
	}
	tx.Commit()
}

func SaveTagItem(tagItem *item_model.TagItem) {
	if (app_repository.Conn.NewRecord(tagItem)) {
		app_repository.Conn.Create(&tagItem)
	} else {
		app_repository.Conn.Save(&tagItem)
	}
}

func TotalUnreadTagItems(tagsIds []string) int {
	count := 0
	app_repository.Conn.Table(tagItemTable).Select("id").Where("unread = ? AND tag_id IN(?)", 1, tagsIds).Count(&count)

	return count
}
