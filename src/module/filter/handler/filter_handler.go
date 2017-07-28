package filter_handler

import (
	"encoding/json"

	"github.com/cstockton/go-conv"

	"github.com/tidinio/src/component/helper/collection"
	"github.com/tidinio/src/component/notifier/fcm"
	"github.com/tidinio/src/module/filter/model"
	"github.com/tidinio/src/module/filter/repository"
)

func AddFilters(userId uint, filtersApi []filter_model.FilterSync, noticeType string) []filter_model.FilterSync {
	newFilters := []filter_model.FilterSync{}
	if len(filtersApi) < 1 {
		return newFilters
	}
	for _, filterApi := range filtersApi {
		filter := addFilter(userId, filterApi)
		filterApi.ID = filter.ID
		newFilters = append(newFilters, filterApi)
	}
	go func() {
		afterAdded(userId, newFilters, noticeType)
	}()
	return newFilters
}

func DeleteFilters(userId uint, filtersApi []filter_model.FilterSync, noticeType string) {
	if len(filtersApi) < 1 {
		return
	}
	filtersDb := filter_repository.GetUserFilters(userId)
	deletedFiltersIds := []uint{}
	modified := false
	for _, filterApi := range filtersApi {
		for _, filterDb := range filtersDb {
			if filterApi.ID == filterDb.ID {
				deleteFilter(filterDb)
				deletedFiltersIds = append(deletedFiltersIds, filterDb.ID)
			}
			modified = true
		}
	}
	if modified {
		go func() {
			afterDeleted(userId, deletedFiltersIds, noticeType)
		}()
	}
}

func GetFilters(userId uint) []filter_model.FilterSync {
	filtersDb := filter_repository.GetUserFilters(userId)
	filterFeedsDb := filter_repository.GetFilterFeedsByUserId(userId)
	return filter_model.JoinToFiltersSync(filtersDb, filterFeedsDb)
}

func UpdateFilters(userId uint, filtersSync []filter_model.FilterSync, noticeType string) {
	if len(filtersSync) < 1 {
		return
	}
	filterFeedsDb := getUserFiltersUpdate(userId)
	updatedFilters := []filter_model.FilterSync{}
	modified := false
	for _, filterSync := range filtersSync {
		for _, filterFeedDb := range filterFeedsDb {
			if filterSync.ID != filterFeedDb.ID {
				continue
			}
			if filterSync.Name != filterFeedDb.Name {
				updateFilter(filterFeedDb, filterSync)
			}
			if !collection_helper.EqualeUintSlice(filterSync.FeedIds, filterFeedDb.FeedIds) {
				updateFilterFeeds(filterSync.ID, filterSync.FeedIds)
			}
			updatedFilters = append(updatedFilters, filterSync)
			modified = true
		}
	}
	if modified {
		go func() {
			afterUpdated(userId, updatedFilters, noticeType)
		}()
	}
}

func addFilter(userId uint, filterApi filter_model.FilterSync) filter_model.Filter {
	filter := filter_model.Filter{}
	filter.Type = filterApi.Type
	filter.Name = filterApi.Name
	filter.Enabled = conv.Int(filterApi.Enabled)
	filter.TagId = filterApi.TagId
	filter.UserId = userId
	filter_repository.SaveFilter(&filter)
	for _, feedId := range filterApi.FeedIds {
		addFilterFeed(filter.ID, feedId)
	}
	return filter
}

func afterAdded(userId uint, filtersApi []filter_model.FilterSync, noticeType string) {
	data, err := json.Marshal(filtersApi)
	if err != nil {
		return
	}
	app_fcm.Send(userId, app_fcm.AddFilters, string(data), noticeType)
}

func afterDeleted(userId uint, filtersApiId []uint, noticeType string) {
	data, err := json.Marshal(filtersApiId)
	if err != nil {
		return
	}
	app_fcm.Send(userId, app_fcm.DeleteFilters, string(data), noticeType)
}

func afterUpdated(userId uint, filtersApi []filter_model.FilterSync, noticeType string) {
	data, err := json.Marshal(filtersApi)
	if err != nil {
		return
	}
	app_fcm.Send(userId, app_fcm.UpdateFilters, string(data), noticeType)
}

func deleteFilter(filter filter_model.Filter) {
	filter_repository.DeleteFilterFeeds(filter.ID)
	filter_repository.DeleteFilter(filter)
}

func getUserFiltersUpdate(userId uint) []filter_model.FilterUpdate {
	filtersDb := filter_repository.GetUserFilters(userId)
	filterFeedsDb := filter_repository.GetFilterFeedsByUserId(userId)
	return filter_model.JoinToFiltersUpdate(filtersDb, filterFeedsDb)
}

func updateFilter(filterUpdate filter_model.FilterUpdate, filterSync filter_model.FilterSync) {
	filter := filter_model.JoinToUpdate(filterUpdate, filterSync)
	filter_repository.SaveFilter(&filter)
}

func updateFilterFeeds(filterId uint, feedsId []uint) {
	filter_repository.DeleteFilterFeeds(filterId)
	for _, feedsId := range feedsId {
		filterFeed := filter_model.FilterFeed{}
		filterFeed.FilterId = filterId
		filterFeed.FeedId = feedsId
		filter_repository.SaveFilterFeed(&filterFeed)
	}
}
