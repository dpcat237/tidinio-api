package feed_repository

import (
	"github.com/tidinio/src/core/feed/model"
	"github.com/tidinio/src/core/component/repository"
)

func GetLastHistory(repo common_repository.Repository, feedId uint) feed_model.FeedHistory {
	feedHistory := feed_model.FeedHistory{}
	repo.DB.Where("feed_id = ?", feedId).First(&feedHistory)

	return feedHistory
}

func SaveFeedHistory(repo common_repository.Repository, feedHistory *feed_model.FeedHistory) {
	if (repo.DB.NewRecord(feedHistory)) {
		repo.DB.Create(&feedHistory)
	} else {
		repo.DB.Save(&feedHistory)
	}
}
