package feed_repository

import (
	"github.com/tidinio/src/module/feed/model"
	"github.com/tidinio/src/component/repository"
)

const userFeedTable  = feed_model.UserFeedTable

func GetUserFeedByFeedAndUser(feedId uint, userId uint) feed_model.UserFeed {
	userFeed := feed_model.UserFeed{}
	app_repository.Conn.Where("feed_id = ? AND user_id = ?", feedId, userId).First(&userFeed)

	return userFeed
}

func GetUserFeedById(userFeedId uint) feed_model.UserFeed {
	userFeed := feed_model.UserFeed{}
	app_repository.Conn.Where("id = ?", userFeedId).First(&userFeed)

	return userFeed
}

func GetUserFeedsByUserId(userId uint) []feed_model.UserFeed {
	userFeeds := []feed_model.UserFeed{}
	app_repository.Conn.Table(userFeedTable).Where("user_id = ?", userId).Scan(&userFeeds)

	return userFeeds
}

func SaveUserFeed(userFeed *feed_model.UserFeed) {
	if (app_repository.Conn.NewRecord(userFeed)) {
		app_repository.Conn.Create(&userFeed)
	} else {
		app_repository.Conn.Save(&userFeed)
	}
}
