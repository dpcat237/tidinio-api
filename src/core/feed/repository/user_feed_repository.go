package feed_repository

import (
	"github.com/tidinio/src/core/feed/model"
	"github.com/tidinio/src/core/component/repository"
)

func GetUserFeedByFeedAndUser(repo common_repository.Repository, feedId uint, userId uint) feed_model.UserFeed {
	userFeed := feed_model.UserFeed{}
	repo.DB.Where("feed_id = ? AND user_id = ?", feedId, userId).First(&userFeed)

	return userFeed
}

func SaveUserFeed(repo common_repository.Repository, userFeed *feed_model.UserFeed) {
	if (repo.DB.NewRecord(userFeed)) {
		repo.DB.Create(&userFeed)
	} else {
		repo.DB.Save(&userFeed)
	}
}
