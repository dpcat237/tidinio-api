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

func SaveTagItem(repo common_repository.Repository, tagItem *item_model.TagItem) {
	if (repo.DB.NewRecord(tagItem)) {
		repo.DB.Create(&tagItem)
	} else {
		repo.DB.Save(&tagItem)
	}
}
