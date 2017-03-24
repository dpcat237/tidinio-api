package item_repository

import (
	"github.com/jinzhu/gorm"
	"github.com/tidinio/src/core/item/model"
	"github.com/tidinio/src/core/component/repository"
	"fmt"
	"time"
)

const userItemTable = item_model.UserItemTable

type UserItemRepository struct {
	db *gorm.DB
}

func CountUnreadByUser(repo UserItemRepository, userId uint) int {
	count := 0
	repo.db.Table("user_item").Select("id").Where("shared = ? AND unread = ? AND user_id = ?", 0, 1, userId).Count(&count)

	return count
}

func Close(userItemRepo UserItemRepository) {
	userItemRepo.db.Close()
}

func GetUnreadUserItems(
userItemRepo UserItemRepository,
userId uint,
unreadIds []string,
offset int,
limit int) []item_model.UserItem {
	results := []item_model.UserItem{}
	iTb := "item"
	ufTb := "user_feed"
	userItemRepo.db.
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

func NewUserItemRepository() UserItemRepository {
	userItemRepo := UserItemRepository{}
	userItemRepo.db = common_repository.InitConnection()

	return userItemRepo
}

func SyncReadItems(repo UserItemRepository, items []item_model.UserItem) {
	tx := repo.db.Begin()
	now := time.Now().Format("2006-01-02 15:04:05")

	for _, item := range items {
		query := fmt.Sprintf(
			"UPDATE user_item SET stared='%d', unread='%d', updated_at='%s' WHERE id='%d';",
			item.Stared, item.Unread, now, item.ID)
		tx.Exec(query)
	}
	tx.Commit()
}
