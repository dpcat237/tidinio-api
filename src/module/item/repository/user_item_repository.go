package item_repository

import (
	"fmt"
	"github.com/cstockton/go-conv"
	"github.com/tidinio/src/module/item/model"
	"github.com/tidinio/src/component/repository"
)

const userItemTable = item_model.UserItemTable

func CountUnreadByUser(userId uint) int {
	count := 0
	app_repository.Conn.Table(userItemTable).Select("count(id)").
		Where("shared = ? AND unread = ? AND user_id = ?", 0, 1, userId).Count(&count)
	return count
}

func GetUnreadUserItems(
	userId uint,
	unreadIds []string,
	offset int,
	limit int) []item_model.UserItem {
	results := []item_model.UserItem{}
	iTb := "item"
	ufTb := "user_feed"
	app_repository.Conn.
		Table(userItemTable).
		Select(userItemTable + ".id, " + userItemTable + ".stared, " + userItemTable + ".unread, " + userItemTable + ".item_id").
		Joins(
		"inner join " + iTb + " on " + userItemTable + ".item_id = " + iTb + ".id "+
			"and "+ userItemTable+ ".shared = ? and "+ userItemTable+ ".unread = ? "+
			"and "+ userItemTable+ ".user_id = ? AND "+ userItemTable+ ".id NOT IN(?) "+
			"left join "+ ufTb+ " on "+ iTb+ ".feed_id="+ ufTb+ ".feed_id "+
			"and "+ ufTb+ ".user_id = ? and "+ ufTb+ ".deleted_at = ?", 0, 1, userId, unreadIds, userId, nil).
		Order("" + userItemTable + ".item_id desc").
		Offset(offset).
		Limit(limit).
		Scan(&results)
	return results
}

func GetUnreadUserItemsContent(unreadIds []string) []item_model.TagItemList {
	results := []item_model.TagItemList{}
	iTb := "item"
	fTb := "feed"
	app_repository.Conn.
		Table(userItemTable).
		Select(userItemTable + ".id user_item_id, " + fTb + ".id feed_id, " + userItemTable + ".stared, " + fTb +
		".language, " + iTb + ".link, " + iTb + ".title, " + iTb + ".content, " + iTb + ".created_at, " + iTb + ".published_at").
		Joins(
		"inner join " + iTb + " on " + userItemTable + ".item_id = " + iTb + ".id "+
			"and "+ userItemTable+ ".id IN(?) "+
			"left join "+ fTb+ " on "+ iTb+ ".feed_id="+ fTb+ ".id ", unreadIds).
		Scan(&results)
	return results
}

func GetUserItemByItemUserId(itemId uint, userId uint) item_model.UserItem {
	userItem := item_model.UserItem{}
	app_repository.Conn.Where("item_id = ? AND user_id = ?", itemId, userId).First(&userItem)

	return userItem
}

func SaveUserItem(userItem *item_model.UserItem) {
	if app_repository.Conn.NewRecord(userItem) {
		app_repository.Conn.Create(&userItem)
	} else {
		app_repository.Conn.Save(&userItem)
	}
}

func SyncReadItems(items []item_model.UserItem) {
	tx := app_repository.Conn.Begin()
	for _, item := range items {
		query := fmt.Sprintf(
			"UPDATE user_item SET stared='%d', unread='%d', updated_at='%s' WHERE id='%d';",
			item.Stared, item.Unread, app_repository.GetDateNowFormatted(), item.ID)
		tx.Exec(query)
	}
	tx.Commit()
}

func UpdateUserItemsStaredStatus(items []item_model.TagItemList) {
	tx := app_repository.Conn.Begin()
	for _, item := range items {
		stared := conv.Int(item.Stared)
		query := fmt.Sprintf(
			"UPDATE user_item SET stared='%d', updated_at='%s' WHERE id='%d';",
			stared, app_repository.GetDateNowFormatted(), item.ArticleId)
		tx.Exec(query)
	}
	tx.Commit()
}
