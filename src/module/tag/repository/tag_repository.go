package tag_repository

import (
	"github.com/tidinio/src/component/repository"
	"github.com/tidinio/src/module/tag/model"
)

const tagTable = tag_model.TagTable

func DeleteTag(tag tag_model.Tag) {
	app_repository.Conn.Delete(tag)
}

func GetUserTags(userId uint) []tag_model.Tag {
	tags := []tag_model.Tag{}
	app_repository.Conn.Table(tagTable).Where("user_id = ? and deleted_at IS NULL", userId).Scan(&tags)

	return tags
}

func SaveTag(tag *tag_model.Tag) {
	if app_repository.Conn.NewRecord(tag) {
		app_repository.Conn.Create(&tag)
	} else {
		app_repository.Conn.Save(&tag)
	}
}
