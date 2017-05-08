package feed_repository

import (
	"github.com/tidinio/src/module/feed/model"
	"github.com/tidinio/src/component/repository"
)

func GetFeedSources() []feed_model.FeedSourceCategory {
	feedSourceCategories := []feed_model.FeedSourceCategory{}
	app_repository.Conn.Preload("FeedSources").Find(&feedSourceCategories).Limit(2)

	return feedSourceCategories
}
