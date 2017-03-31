package feed_handler

import (
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/feed/repository"
	"github.com/tidinio/src/core/feed/model"
)

func SubscribeUserToFeed(userId uint, feed feed_model.Feed) {
	repo := common_repository.InitConnection()
	subscribeNewUser(userId, feed.ID)
	if (!feed.IsEnabled()) {
		feed.Enable()
		feed_repository.SaveFeed(repo, &feed)
	}

	defer repo.Close()
}

func subscribeNewUser(userId uint, feedId uint) {
	repo := common_repository.InitConnection()
	feed := feed_repository.GetFeedById(repo, feedId)
	userFeed := feed_model.UserFeed{}
	userFeed.UserId = userId
	userFeed.FeedId = feedId
	userFeed.Title = feed.Title
	feed_repository.SaveUserFeed(repo, &userFeed)
}
