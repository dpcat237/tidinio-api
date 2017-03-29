package item_repository

import (
	"github.com/tidinio/src/core/item/model"
	"github.com/tidinio/src/core/component/repository"
	"time"
	"fmt"
)

const tagItemTable  = item_model.TagItemTable

func CreateTagItems(repo common_repository.Repository, collection []item_model.TagItem) {
	tx := repo.DB.Begin()
	for _, tagItem := range collection {
		tx.Create(&tagItem)
	}
	tx.Commit()
}

func GetTagItemByUserItemTagId(repo common_repository.Repository, userItemId uint, tagId uint) item_model.TagItem {
	tagItem := item_model.TagItem{}
	repo.DB.Where("user_item_id = ? AND later_id = ?", userItemId, tagId).First(&tagItem)

	return tagItem
}

func GetTagsByUserItemIds(repo common_repository.Repository, userItemIds []string, unread int) []item_model.TagItem {
	results := []item_model.TagItem{}
	repo.DB.
		Table(tagItemTable).
		Select("id, user_item_id, later_id").
		Where("unread = ? AND user_item_id IN(?)", unread, userItemIds).
		Scan(&results)

	return results
}

func GetUnreadTagItems(
repo common_repository.Repository,
tagsIds []string,
tagItemsIds []string,
offset int,
limit int) []item_model.TagItem {
	results := []item_model.TagItem{}
	repo.DB.
		Table(tagItemTable).
		Select("id, user_item_id, unread").
		Where("unread = ? AND later_id IN(?) AND id NOT IN(?)", 1, tagsIds, tagItemsIds).
		Offset(offset).
		Limit(limit).
		Scan(&results)

	return results
}

func GetUnreadTagItemsSync(repo common_repository.Repository, userItemIds []string) []item_model.TagItemSyncDB {
	results := []item_model.TagItemSyncDB{}
	repo.DB.
		Table(tagItemTable).
		Select(tagItemTable + ".id, " + tagItemTable + ".user_item_id, " + tagItemTable + ".later_id, " + userItemTable +
		".stared").
		Joins(
		"inner join " + userItemTable + " on " + tagItemTable + ".user_item_id = " + userItemTable + ".id " +
			"AND " + userItemTable + ".id IN(?) AND " + tagItemTable + ".unread = ?", userItemIds, 1).
		Scan(&results)

	return results
}

func MarkAsUnread(repo common_repository.Repository, collection map[uint][]uint, unread int) {
	tx := repo.DB.Begin()
	now := time.Now().Format("2006-01-02 15:04:05")

	for userItemId, tagsId := range collection {
		for _, tagId := range tagsId {
			query := fmt.Sprintf(
				"UPDATE later_item SET unread='%d', updated_at='%s' WHERE user_item_id='%d' AND later_id='%d';",
				unread, now, userItemId, tagId)
			tx.Exec(query)
		}
	}
	tx.Commit()
}

func SaveTagItem(repo common_repository.Repository, tagItem *item_model.TagItem) {
	if (repo.DB.NewRecord(tagItem)) {
		repo.DB.Create(&tagItem)
	} else {
		repo.DB.Save(&tagItem)
	}
}

func TotalUnreadTagItems(repo common_repository.Repository, tagsIds []string) int {
	count := 0
	repo.DB.Table(tagItemTable).Select("id").Where("unread = ? AND later_id IN(?)", 1, tagsIds).Count(&count)

	return count
}
