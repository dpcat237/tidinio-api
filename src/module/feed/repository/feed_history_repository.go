package feed_repository

import (
	"github.com/tidinio/src/module/feed/model"
	"github.com/tidinio/src/component/repository"
)

func GetLastHistory(feedId uint) feed_model.FeedHistory {
	feedHistory := feed_model.FeedHistory{}
	app_repository.Conn.Where("feed_id = ?", feedId).First(&feedHistory)

	return feedHistory
}

func SaveFeedHistory(feedHistory *feed_model.FeedHistory) {
	if (app_repository.Conn.NewRecord(feedHistory)) {
		app_repository.Conn.Create(&feedHistory)
	} else {
		app_repository.Conn.Save(&feedHistory)
	}
}
