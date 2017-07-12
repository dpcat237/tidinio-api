package feed_repository

import (
	"github.com/tidinio/src/module/feed/model"
	"github.com/tidinio/src/component/repository"
)

func GetFeedById(feedId uint) feed_model.Feed {
	feed := feed_model.Feed{}
	app_repository.Conn.Where("id = ?", feedId).First(&feed)

	return feed
}

func GetFeedByUrl(link string) feed_model.Feed {
	feed := feed_model.Feed{}
	app_repository.Conn.Where("url = ?", link).First(&feed)

	return feed
}

func SaveFeed(feed *feed_model.Feed) {
	if app_repository.Conn.NewRecord(feed) {
		app_repository.Conn.Create(&feed)
	} else {
		app_repository.Conn.Save(&feed)
	}
}
