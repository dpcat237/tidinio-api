package feed_repository

import (
	"github.com/tidinio/src/core/feed/model"
	"github.com/tidinio/src/core/component/repository"
)

func GetUserFeedByFeedAndUser(repo app_repository.Repository, feedId uint, userId uint) feed_model.UserFeed {
	userFeed := feed_model.UserFeed{}
	repo.DB.Where("feed_id = ? AND user_id = ?", feedId, userId).First(&userFeed)

	return userFeed
}

func GetUserFeedById(repo app_repository.Repository, userFeedId uint) feed_model.UserFeed {
	userFeed := feed_model.UserFeed{}
	repo.DB.Where("id = ?", userFeedId).First(&userFeed)

	return userFeed
}

func GetUserFeedsByUserId(repo app_repository.Repository, userId uint) []feed_model.UserFeed {
	userFeeds := []feed_model.UserFeed{}
	repo.DB.Where("user_id = ?", userId).Scan(&userFeeds)

	return userFeeds
}

func SaveUserFeed(repo app_repository.Repository, userFeed *feed_model.UserFeed) {
	if (repo.DB.NewRecord(userFeed)) {
		repo.DB.Create(&userFeed)
	} else {
		repo.DB.Save(&userFeed)
	}
}
