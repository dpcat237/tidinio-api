package feed_handler

import (
	"github.com/tidinio/src/module/feed/repository"
	"github.com/tidinio/src/module/feed/model"
)

func GetFeedSources() []feed_model.FeedSourceCategorySync {
	return feed_model.ToFeedSourceCategorySyncs(feed_repository.GetFeedSources())
}
