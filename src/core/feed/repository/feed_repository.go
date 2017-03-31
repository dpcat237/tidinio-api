package feed_repository

import (
	"github.com/tidinio/src/core/feed/model"
	"github.com/tidinio/src/core/component/repository"
)

func GetFeedById(repo common_repository.Repository, feedId uint) feed_model.Feed {
	feed := feed_model.Feed{}
	repo.DB.Where("id = ?", feedId).First(&feed)

	return feed
}

func GetFeedByUrl(repo common_repository.Repository, link string) feed_model.Feed {
	feed := feed_model.Feed{}
	repo.DB.Where("url = ?", link).First(&feed)

	return feed
}

func SaveFeed(repo common_repository.Repository, feed *feed_model.Feed) {
	if (repo.DB.NewRecord(feed)) {
		repo.DB.Create(&feed)
	} else {
		repo.DB.Save(&feed)
	}
}
