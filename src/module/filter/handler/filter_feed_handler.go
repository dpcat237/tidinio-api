package filter_handler

import (
	"github.com/tidinio/src/module/filter/model"
	"github.com/tidinio/src/module/filter/repository"
)

func addFilterFeed(filterId uint, feedId uint) filter_model.FilterFeed {
	filterFeed := filter_model.FilterFeed{}
	filterFeed.FilterId = filterId
	filterFeed.FeedId = feedId
	filter_repository.SaveFilterFeed(&filterFeed)
	return filterFeed
}
