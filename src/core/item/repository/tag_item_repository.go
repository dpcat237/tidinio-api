package item_repository

import (
	"github.com/tidinio/src/core/item/model"
	"github.com/tidinio/src/core/component/repository"
)

const tagItemTable  = item_model.TagItemTable

func GetTagItemByUserItemTagId(repo common_repository.Repository, userItemId uint, tagId uint) item_model.TagItem {
	tagItem := item_model.TagItem{}
	repo.DB.Where("user_item_id = ? AND later_id = ?", userItemId, tagId).First(&tagItem)

	return tagItem
}

func GetTagsByUserItemIds(repo common_repository.Repository, userItemIds []string) []item_model.TagItem {
	results := []item_model.TagItem{}
	repo.DB.
		Table(tagItemTable).
		Select("id, user_item_id, later_id").
		Where("unread = ? AND user_item_id IN(?)", 1, userItemIds).
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
