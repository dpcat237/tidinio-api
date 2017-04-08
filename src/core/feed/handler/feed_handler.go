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
	"github.com/tidinio/src/core/item/repository"
	"github.com/tidinio/src/core/component/helper/http"
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

	return feed, feedError
}

func EditFeedTitle(userId uint, userFeedId uint, feedTitle string) error {
	repo := app_repository.InitConnection()
	userFeed := feed_repository.GetUserFeedById(repo, userFeedId)
	if (userFeed.ID < 1 || userFeed.UserId != userId) {
		return errors.New("Wrong provided data")
	}

	if (userFeed.Title == feedTitle) {
		return nil
	}

	userFeed.Title = feedTitle
	feed_repository.SaveUserFeed(repo, &userFeed)
	go func() {
		afterUserFeedModified(userId)
	}()

	return nil
}

func afterFeedCreated(feed feed_model.Feed) {
	feed.Language = detectFeedLanguage(feed.ID)
	crawling := detectFeedNeedCrawling(feed.ID)
	if (crawling) {
		feed.Crawling = app_repository.BoolToInt(crawling)
	}

	repo := app_repository.InitConnection()
	feed_repository.SaveFeed(repo, &feed)
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

func detectFeedLanguage(feedId uint) string {
	language := ""
	repo := app_repository.InitConnection()
	items := item_repository.GetLastItems(repo, feedId, 10)
	defer repo.Close()
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
	repo := app_repository.InitConnection()
	items := item_repository.GetLastItems(repo, feedId, 5)
	defer repo.Close()

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
