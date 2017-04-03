package feed_handler

import (
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/feed/repository"
	"github.com/tidinio/src/core/feed/model"
	"time"
)

func addHistory(repo app_repository.Repository, feedId uint) {
	feedHistory := feed_model.FeedHistory{}
	feedHistory.FeedId = feedId
	feed_repository.SaveFeedHistory(repo, &feedHistory)
}

func updateHistoryDataChanged(feedId uint) {
	repo := app_repository.InitConnection()
	feedHistory := feed_repository.GetLastHistory(repo, feedId)
	if (feedHistory.ID > 0) {
		feedHistory.SetFinished()
		feed_repository.SaveFeedHistory(repo, &feedHistory)
	}
	defer repo.Close()
}

func updateHistorySameData(feedId uint) {
	repo := app_repository.InitConnection()
	feedHistory := feed_repository.GetLastHistory(repo, feedId)
	if (feedHistory.ID < 1) {
		addHistory(repo, feedId)
		defer repo.Close()

		return
	}

	historyLimit := time.Now().AddDate(0, 0, 1)
	if (feedHistory.IsFinished() || feedHistory.CreatedAt.Before(historyLimit)) {
		addHistory(repo, feedId)
		defer repo.Close()

		return
	}

	feedHistory.IncreaseCounter()
	feed_repository.SaveFeedHistory(repo, &feedHistory)
	defer repo.Close()
}
