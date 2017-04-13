package feed_handler

import (
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/feed/repository"
	"github.com/tidinio/src/core/feed/model"
	"errors"
	"github.com/tidinio/src/core/feed/data_transformer"
	"github.com/tidinio/src/core/component/notifier/fcm"
)

func SyncUserFeeds(userId uint, userFeedsApi []feed_model.UserFeedSync) []feed_model.UserFeedSync {
	repo := app_repository.InitConnection()
	userFeeds := feed_repository.GetUserFeedsByUserId(repo, userId)
	if len(userFeeds) < 1 {
		return []feed_model.UserFeedSync{}
	}

	if len(userFeedsApi) < 1 {
		return feed_transformer.ToUserFeedsSync(userFeeds)
	}

	modified := false
	for _, userFeed := range userFeeds {
		for _, userFeedApi := range userFeedsApi {
			if (userFeed.ID != userFeedApi.ID) {
				continue
			}

			if (userFeed.Title != userFeedApi.Title && userFeed.UpdatedAt.Before(userFeedApi.UpdatedAt)) {
				userFeed.Title = userFeedApi.Title
				feed_repository.SaveUserFeed(repo, &userFeed)
				modified = true
			}
		}
	}

	if (modified) {
		go func() {
			afterUserFeedModified(userId)
		}()
	}

	return feed_transformer.ToUserFeedsSync(feed_repository.GetUserFeedsByUserId(repo, userId))
}

func SubscribeUserToFeed(userId uint, feed feed_model.Feed) {
	repo := app_repository.InitConnection()
	subscribeNewUser(userId, feed.ID)
	if (!feed.IsEnabled()) {
		feed.Enable()
		feed_repository.SaveFeed(repo, &feed)
	}

	defer repo.Close()
}

func UnsubscribeFromFeed(userId uint, userFeedId uint) error {
	repo := app_repository.InitConnection()
	userFeed := feed_repository.GetUserFeedById(repo, userFeedId)
	if (userFeed.ID < 1 || userFeed.UserId != userId) {
		return errors.New("Wrong provided data")
	}

	userFeed.DeletedAt = app_repository.GetDateNow()
	feed_repository.SaveUserFeed(repo, &userFeed)
	go func() {
		afterUserFeedModified(userId)
	}()

	return nil
}

func afterUserFeedModified(userId uint) {
	app_fcm.RequireToSync(app_fcm.SyncFeeds, userId)
}

func subscribeNewUser(userId uint, feedId uint) {
	repo := app_repository.InitConnection()
	feed := feed_repository.GetFeedById(repo, feedId)
	userFeed := feed_model.UserFeed{}
	userFeed.UserId = userId
	userFeed.FeedId = feedId
	userFeed.Title = feed.Title
	feed_repository.SaveUserFeed(repo, &userFeed)
}
