package tag_handler

import (
	"encoding/json"

	"github.com/tidinio/src/component/notifier/fcm"
	"github.com/tidinio/src/module/tag/model"
	"github.com/tidinio/src/module/tag/repository"
)

func AddTags(userId uint, tagsApi []tag_model.TagSync, noticeType string) []tag_model.TagSync {
	newTags := []tag_model.TagSync{}
	if len(tagsApi) < 1 {
		return newTags
	}

	for _, tagApi := range tagsApi {
		newTags = append(newTags, addTag(userId, tagApi.Name).ConvertToTagSync())
	}
	go func() {
		afterAdded(userId, newTags, noticeType)
	}()
	return newTags
}

func DeleteTags(userId uint, tagsApi []tag_model.TagSync, noticeType string) {
	if len(tagsApi) < 1 {
		return
	}
	tagsDb := tag_repository.GetUserTags(userId)
	deletedTagsIds := []uint{}
	modified := false
	for _, tagApi := range tagsApi {
		for _, tagDb := range tagsDb {
			if tagApi.ID == tagDb.ID {
				deleteTag(tagDb)
				deletedTagsIds = append(deletedTagsIds, tagDb.ID)
			}
			modified = true
		}
	}
	if modified {
		go func() {
			afterDeleted(userId, deletedTagsIds, noticeType)
		}()
	}
}

func GetTags(userId uint) []tag_model.TagSync {
	tagsDb := tag_repository.GetUserTags(userId)
	return tag_model.ConvertToTagsSync(tagsDb)
}

func UpdateTags(userId uint, tagsApi []tag_model.TagSync, noticeType string) {
	if len(tagsApi) < 1 {
		return
	}
	tagsDb := tag_repository.GetUserTags(userId)
	updatedTags := []tag_model.TagSync{}
	modified := false
	for _, tagApi := range tagsApi {
		for _, tagDb := range tagsDb {
			if tagApi.ID != tagDb.ID {
				continue
			}
			if tagApi.Name != tagDb.Name && tagDb.UpdatedAt.Before(tagApi.UpdatedAt) {
				updatedTags = append(updatedTags, updateTag(tagDb, tagApi.Name).ConvertToTagSync())
				modified = true
			}
		}
	}
	if modified {
		go func() {
			afterUpdated(userId, updatedTags, noticeType)
		}()
	}
}

func addTag(userId uint, name string) tag_model.Tag {
	tag := tag_model.Tag{}
	tag.Name = name
	tag.UserId = userId
	tag.Enabled = 1
	tag_repository.SaveTag(&tag)
	return tag
}

func afterAdded(userId uint, tagsApi []tag_model.TagSync, noticeType string) {
	data, err := json.Marshal(tagsApi)
	if err != nil {
		return
	}
	app_fcm.Send(userId, app_fcm.AddTags, string(data), noticeType)
}

func afterDeleted(userId uint, tagsApiId []uint, noticeType string) {
	data, err := json.Marshal(tagsApiId)
	if err != nil {
		return
	}
	app_fcm.Send(userId, app_fcm.DeleteTags, string(data), noticeType)
}

func afterUpdated(userId uint, tagsApi []tag_model.TagSync, noticeType string) {
	data, err := json.Marshal(tagsApi)
	if err != nil {
		return
	}
	app_fcm.Send(userId, app_fcm.UpdateTags, string(data), noticeType)
}

func deleteTag(tag tag_model.Tag) {
	tag_repository.DeleteTag(tag)
}

func updateTag(tag tag_model.Tag, newName string) tag_model.Tag {
	tag.Name = newName
	tag_repository.SaveTag(&tag)
	return tag
}
