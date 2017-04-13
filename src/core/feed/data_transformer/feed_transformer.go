package feed_transformer

import "github.com/tidinio/src/core/feed/model"

func ToUserFeedsSync(userFeeds []feed_model.UserFeed) []feed_model.UserFeedSync {
	userFeedsSync := []feed_model.UserFeedSync{}
	for _, userFeed := range userFeeds {
		userFeedsSync = append(userFeedsSync, ToUserFeedSync(userFeed))
	}

	return userFeedsSync
}

func ToUserFeedSync(userFeed feed_model.UserFeed) feed_model.UserFeedSync {
	userFeedSync := feed_model.UserFeedSync{}
	userFeedSync.ID = userFeed.ID
	userFeedSync.Title = userFeed.Title
	userFeedSync.UpdatedAt = userFeed.UpdatedAt

	return userFeedSync
}
