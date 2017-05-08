package feed_handler

import (
	"errors"
	"github.com/tidinio/src/module/feed/model"
	"github.com/tidinio/src/module/feed/repository"
	"github.com/tidinio/src/component/repository"
	"github.com/tidinio/src/component/notifier/fcm"
)

func SyncUserFeeds(userId uint, userFeedsApi []feed_model.UserFeedSync) []feed_model.UserFeedSync {
	userFeeds := feed_repository.GetUserFeedsByUserId(userId)
	if len(userFeeds) < 1 {
		return []feed_model.UserFeedSync{}
	}

	if len(userFeedsApi) < 1 {
		return feed_model.ToUserFeedsSync(userFeeds)
	}

	modified := false
	for _, userFeed := range userFeeds {
		for _, userFeedApi := range userFeedsApi {
			if (userFeed.ID != userFeedApi.ID) {
				continue
			}

			if (userFeed.Title != userFeedApi.Title && userFeed.UpdatedAt.Before(userFeedApi.UpdatedAt)) {
				userFeed.Title = userFeedApi.Title
				feed_repository.SaveUserFeed(&userFeed)
				modified = true
			}
		}
	}

	if (modified) {
		go func() {
			afterUserFeedModified(userId)
		}()
	}

	return feed_model.ToUserFeedsSync(feed_repository.GetUserFeedsByUserId(userId))
}

func SubscribeUserToFeed(userId uint, feed feed_model.Feed) {
	subscribeNewUser(userId, feed.ID)
	if (!feed.IsEnabled()) {
		feed.Enable()
		feed_repository.SaveFeed(&feed)
	}
}

func UnsubscribeFromFeed(userId uint, userFeedId uint) error {
	userFeed := feed_repository.GetUserFeedById(userFeedId)
	if (userFeed.ID < 1 || userFeed.UserId != userId) {
		return errors.New("Wrong provided data")
	}

	userFeed.DeletedAt = app_repository.GetDateNow()
	feed_repository.SaveUserFeed(&userFeed)
	go func() {
		afterUserFeedModified(userId)
	}()

	return nil
}

func afterUserFeedModified(userId uint) {
	app_fcm.RequireToSync(app_fcm.SyncFeeds, userId)
}

func subscribeNewUser(userId uint, feedId uint) {
	feed := feed_repository.GetFeedById(feedId)
	userFeed := feed_model.UserFeed{}
	userFeed.UserId = userId
	userFeed.FeedId = feedId
	userFeed.Title = feed.Title
	feed_repository.SaveUserFeed(&userFeed)
}
