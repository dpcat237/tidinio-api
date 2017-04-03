package feed_handler

import (
	"errors"
	"github.com/mmcdole/gofeed"
	"github.com/tidinio/src/core/feed/model"
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/feed/repository"
	"github.com/tidinio/src/core/item/handler"
	"github.com/tidinio/src/core/feed/redis"
	"github.com/tidinio/src/core/component/helper/string"
)

func AddFeed(userId uint, feedUrl string) (feed_model.Feed, error) {
	feed := feed_model.Feed{}
	if (feedUrl == "") {
		return feed, errors.New("Empty url")
	}
	feedUrl, feedError := string_helper.CleanUrl(feedUrl)
	if feedError != nil {
		return feed, feedError
	}

	repo := app_repository.InitConnection()
	feed = feed_repository.GetFeedByUrl(repo, feedUrl)
	if (feed.ID > 0) {
		userFeed := feed_repository.GetUserFeedByFeedAndUser(repo, feed.ID, userId)
		if (userFeed.ID > 0) {
			return feed, nil
		}
	} else if (feed.ID < 1) {
		feed, feedError = createFeed(repo, feedUrl)
		updateFeedData(feed)
		//TODO: NPSCoreEvents::FEED_CREATED
	}

	isFeedEnabled := feed.IsEnabled()
	SubscribeUserToFeed(userId, feed)
	if (!isFeedEnabled) {
		updateFeedData(feed)
	}
	item_handler.AddLastItemsNewUser(userId, feed.ID, 25)
	//TODO: NPSCoreEvents::FEED_MODIFIED

	return feed, feedError
}

func createFeed(repo app_repository.Repository, feedUrl string) (feed_model.Feed, error) {
	feed := feed_model.Feed{}
	fp := gofeed.NewParser()
	feedData, feedError := fp.ParseURL(feedUrl)
	if feedError != nil {
		return feed, errors.New("Feed with wrong data")
	}

	feed.Url = feedUrl
	feed.Title = feedData.Title
	feed.Website = feedData.Link
	feed.Enable()
	feed_repository.SaveFeed(repo, &feed)

	return feed, nil
}

func isFeedDataChanged(feed feed_model.Feed) bool {
	currentHash := string_helper.GetHashFromUrlData(feed.Url)
	lastHash := feed_redis.GetDataHash(feed.ID)
	if (currentHash == lastHash) {
		updateHistorySameData(feed.ID)

		return false
	}
	feed_redis.SetDataHash(feed.ID, currentHash)

	return true
}

func updateFeedData(feed feed_model.Feed) {
	if (!isFeedDataChanged(feed)) {
		return
	}

	fp := gofeed.NewParser()
	feedData, feedError := fp.ParseURL(feed.Url)
	if feedError != nil {
		return
	}

	repo := app_repository.InitConnection()
	count := 0
	for _, item := range feedData.Items {
		isNew := item_handler.CreateUpdateItem(repo, item, feed.ID)
		if (isNew) {
			count++
		}
	}

	feed.UpdateDateSync()
	feed_repository.SaveFeed(repo, &feed)
	if (count > 0) {
		updateHistoryDataChanged(feed.ID)
	} else {
		updateHistorySameData(feed.ID)
	}

	defer repo.Close()
}
