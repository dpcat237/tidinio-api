package feed_handler

import (
	"errors"

	"github.com/cstockton/go-conv"
	"github.com/mmcdole/gofeed"

	"github.com/tidinio/src/component/helper/http"
	"github.com/tidinio/src/component/helper/string"
	"github.com/tidinio/src/module/feed/cache"
	"github.com/tidinio/src/module/feed/model"
	"github.com/tidinio/src/module/feed/repository"
	"github.com/tidinio/src/module/item/handler"
	"github.com/tidinio/src/module/item/repository"
)

func AddFeed(userId uint, feedUrl string) (feed_model.UserFeedSync, error) {
	userFeedSync := feed_model.UserFeedSync{}
	if (feedUrl == "") {
		return userFeedSync, errors.New("Empty url")
	}
	feedUrl, feedError := string_helper.CleanUrl(feedUrl)
	if feedError != nil {
		return userFeedSync, feedError
	}

	feed := feed_repository.GetFeedByUrl(feedUrl)
	if (feed.ID > 0) {
		userFeed := feed_repository.GetUserFeedByFeedAndUser(feed.ID, userId)
		if (userFeed.ID > 0) {
			return feed_model.ToUserFeedSync(userFeed), nil
		}
	} else if (feed.ID < 1) {
		feed, feedError = createFeed(feedUrl)
		updateFeedData(feed)
		go func() {
			afterFeedCreated(feed)
		}()
	}

	isFeedEnabled := feed.IsEnabled()
	SubscribeUserToFeed(userId, feed)
	if (!isFeedEnabled) {
		updateFeedData(feed)
	}
	item_handler.AddLastItemsNewUser(userId, feed.ID, 25)
	go func() {
		afterUserFeedModified(userId)
	}()
	userFeed := feed_repository.GetUserFeedByFeedAndUser(feed.ID, userId)

	return feed_model.ToUserFeedSync(userFeed), feedError
}

func EditFeedTitle(userId uint, userFeedId uint, feedTitle string) error {
	userFeed := feed_repository.GetUserFeedById(userFeedId)
	if (userFeed.ID < 1 || userFeed.UserId != userId) {
		return errors.New("Wrong provided data")
	}

	if (userFeed.Title == feedTitle) {
		return nil
	}

	userFeed.Title = feedTitle
	feed_repository.SaveUserFeed(&userFeed)
	go func() {
		afterUserFeedModified(userId)
	}()

	return nil
}

func afterFeedCreated(feed feed_model.Feed) {
	feed.Language = detectFeedLanguage(feed.ID)
	crawling := detectFeedNeedCrawling(feed.ID)
	if (crawling) {
		feed.Crawling = conv.Int(crawling)
	}

	feed_repository.SaveFeed(&feed)
}

func createFeed(feedUrl string) (feed_model.Feed, error) {
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
	feed_repository.SaveFeed(&feed)

	return feed, nil
}

func detectFeedLanguage(feedId uint) string {
	language := ""
	items := item_repository.GetLastItems(feedId, 10)
	for _, item := range items {
		language = string_helper.DetectLanguageFromHtml(item.Content)
		if (language != "") {
			return language
		}
	}

	return language
}

func detectFeedNeedCrawling(feedId uint) bool {
	needsCrawling := true
	items := item_repository.GetLastItems(feedId, 5)

	for _, item := range items {
		needsCrawling = item_handler.IsItemNeedsCrawling(item)
		if (needsCrawling) {
			return needsCrawling
		}
	}

	return needsCrawling
}

func isFeedDataChanged(feed feed_model.Feed) bool {
	currentHash := http_helper.GetHashFromUrlData(feed.Url)
	lastHash := feed_cache.GetDataHash(feed.ID)
	if (currentHash == lastHash) {
		updateHistorySameData(feed.ID)

		return false
	}
	feed_cache.SetDataHash(feed.ID, currentHash)

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

	count := 0
	for _, item := range feedData.Items {
		isNew := item_handler.CreateUpdateItem(item, feed.ID)
		if (isNew) {
			count++
		}
	}

	feed.UpdateDateSync()
	feed_repository.SaveFeed(&feed)
	if (count > 0) {
		updateHistoryDataChanged(feed.ID)
	} else {
		updateHistorySameData(feed.ID)
	}
}
