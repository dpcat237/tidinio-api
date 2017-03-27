package item_repository

import (
	"github.com/tidinio/src/core/item/model"
	"github.com/tidinio/src/core/component/repository"
	"fmt"
	"time"
)

const userItemTable = item_model.UserItemTable

func CountUnreadByUser(repo common_repository.Repository, userId uint) int {
	count := 0
	repo.DB.Table(userItemTable).Select("id").Where("shared = ? AND unread = ? AND user_id = ?", 0, 1, userId).Count(&count)

	return count
}

func Close(repo common_repository.Repository) {
	repo.DB.Close()
}

func GetUnreadUserItems(
repo common_repository.Repository,
userId uint,
unreadIds []string,
offset int,
limit int) []item_model.UserItem {
	results := []item_model.UserItem{}
	iTb := "item"
	ufTb := "user_feed"
	repo.DB.
		Table(userItemTable).
		Select(userItemTable + ".id, " + userItemTable + ".stared, " + userItemTable + ".unread, " + userItemTable + ".item_id").
		Joins(
		"left join " + iTb + " on " + userItemTable + ".item_id = " + iTb + ".id " +
			"and " + userItemTable + ".shared = ? and " + userItemTable + ".unread = ? " +
			"and " + userItemTable + ".user_id = ? AND " + userItemTable + ".id NOT IN(?) " +
			"left join " + ufTb + " on " + iTb + ".feed_id=" + ufTb + ".feed_id " +
			"and " + userItemTable + ".user_id=" + ufTb + ".user_id and " + ufTb + ".deleted = ?", 0, 1, userId, unreadIds, 0).
		Order("" + userItemTable + ".item_id desc").
		Offset(offset).
		Limit(limit).
		Scan(&results)

	return results
}

func GetUnreadUserItemsContent(repo common_repository.Repository, unreadIds []string) []item_model.TagItemList {
	results := []item_model.TagItemList{}
	iTb := "item"
	fTb := "feed"
	repo.DB.
		Table(userItemTable).
		Select(userItemTable + ".id article_id, " + fTb + ".id feed_id, " + userItemTable + ".stared, " + fTb +
			".language, " + iTb + ".link, " + iTb + ".title, " + iTb + ".content, " + iTb + ".created_at").
		Joins(
			"inner join " + iTb + " on " + userItemTable + ".item_id = " + iTb + ".id " +
			"and " + userItemTable + ".id IN(?) " +
			"left join " + fTb + " on " + iTb + ".feed_id=" + fTb + ".id ", unreadIds).
		Scan(&results)

	return results
}

func GetUserItemByItemUserId(repo common_repository.Repository, itemId uint, userId uint) item_model.UserItem {
	userItem := item_model.UserItem{}
	repo.DB.Where("item_id = ? AND user_id = ?", itemId, userId).First(&userItem)

	return userItem
}

func SaveUserItem(repo common_repository.Repository, userItem *item_model.UserItem) {
	if (repo.DB.NewRecord(userItem)) {
		repo.DB.Create(&userItem)
	} else {
		repo.DB.Save(&userItem)
	}
}

func SyncReadItems(repo common_repository.Repository, items []item_model.UserItem) {
	tx := repo.DB.Begin()
	now := time.Now().Format("2006-01-02 15:04:05")

	for _, item := range items {
		query := fmt.Sprintf(
			"UPDATE user_item SET stared='%d', unread='%d', updated_at='%s' WHERE id='%d';",
			item.Stared, item.Unread, now, item.ID)
		tx.Exec(query)
	}
	tx.Commit()
}
