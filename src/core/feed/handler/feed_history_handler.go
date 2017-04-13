package feed_handler

import (
	"github.com/tidinio/src/core/feed/repository"
	"github.com/tidinio/src/core/feed/model"
	"time"
)

func addHistory(feedId uint) {
	feedHistory := feed_model.FeedHistory{}
	feedHistory.FeedId = feedId
	feed_repository.SaveFeedHistory(&feedHistory)
}

func updateHistoryDataChanged(feedId uint) {
	feedHistory := feed_repository.GetLastHistory(feedId)
	if (feedHistory.ID > 0) {
		feedHistory.SetFinished()
		feed_repository.SaveFeedHistory(&feedHistory)
	}
}

func updateHistorySameData(feedId uint) {
	feedHistory := feed_repository.GetLastHistory(feedId)
	if (feedHistory.ID < 1) {
		addHistory(feedId)

		return
	}

	historyLimit := time.Now().AddDate(0, 0, 1)
	if (feedHistory.IsFinished() || feedHistory.CreatedAt.Before(historyLimit)) {
		addHistory(feedId)

		return
	}

	feedHistory.IncreaseCounter()
	feed_repository.SaveFeedHistory(&feedHistory)
}
