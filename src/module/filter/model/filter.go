package filter_model

import (
	"time"

	"github.com/cstockton/go-conv"
)

const (
	FilterTable = "filter"

	FeedToTagType = "feed_to_tag"
)

type Filter struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name    string
	Type    string
	Enabled int
	UserId  uint `gorm:"column:user_id"`
	TagId   uint `gorm:"column:tag_id"`
}

type FilterSync struct {
	ID        uint      `json:"filter_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Enabled   bool      `json:"enabled"`
	TagId     uint      `json:"tag_id"`
	FeedIds   []uint    `json:"feeds_id"`
	UpdatedAt time.Time `json:"date_up"`
}

type FilterUpdate struct {
	Filter
	FeedIds []uint
}

func (filter Filter) ConvertToFilterSync() FilterSync {
	filterSync := FilterSync{}
	filterSync.ID = filter.ID
	filterSync.Name = filter.Name
	filterSync.Type = filter.Type
	filterSync.Enabled = conv.Bool(filter.Enabled)
	filterSync.TagId = filter.TagId
	filterSync.UpdatedAt = filter.UpdatedAt
	return filterSync
}

func (filter Filter) ConvertToFilterUpdate() FilterUpdate {
	filterUpdate := FilterUpdate{}
	filterUpdate.ID = filter.ID
	filterUpdate.CreatedAt = filter.CreatedAt
	filterUpdate.UpdatedAt = filter.UpdatedAt
	filterUpdate.Name = filter.Name
	filterUpdate.Type = filter.Type
	filterUpdate.Enabled = filter.Enabled
	filterUpdate.UserId = filter.UserId
	filterUpdate.TagId = filter.TagId
	return filterUpdate
}

func JoinToFiltersSync(filters []Filter, filterFeeds []FilterFeed) []FilterSync {
	filtersSync := []FilterSync{}
	for _, filter := range filters {
		filtersSync = append(filtersSync, filter.ConvertToFilterSync())
	}
	return joinFiltersSyncFilterFeeds(filtersSync, filterFeeds)
}

func JoinToFiltersUpdate(filters []Filter, filterFeeds []FilterFeed) []FilterUpdate {
	filtersSync := []FilterUpdate{}
	for _, filter := range filters {
		filtersSync = append(filtersSync, filter.ConvertToFilterUpdate())
	}
	return joinFiltersUpdateFilterFeeds(filtersSync, filterFeeds)
}

func JoinToUpdate(filterUpdate FilterUpdate, filterSync FilterSync) Filter {
	filter := Filter{}
	filter.ID = filterUpdate.ID
	filter.CreatedAt = filterUpdate.CreatedAt
	filter.UpdatedAt = filterSync.UpdatedAt
	filter.Name = filterSync.Name
	filter.Type = filterUpdate.Type
	filter.Enabled = conv.Int(filterSync.Enabled)
	filter.UserId = filterUpdate.UserId
	filter.TagId = filterSync.TagId
	return filter
}

func (Filter) TableName() string {
	return FilterTable
}

func joinFiltersSyncFilterFeeds(filtersSync []FilterSync, filterFeeds []FilterFeed) []FilterSync {
	for _, filterFeed := range filterFeeds {
		for fsKey, filterSync := range filtersSync {
			if filterSync.ID == filterFeed.FilterId {
				filtersSync[fsKey].FeedIds = append(filtersSync[fsKey].FeedIds, filterFeed.FeedId)
			}
		}
	}
	return filtersSync
}

func joinFiltersUpdateFilterFeeds(filtersUpdate []FilterUpdate, filterFeeds []FilterFeed) []FilterUpdate {
	for _, filterFeed := range filterFeeds {
		for fsKey, filterUpdate := range filtersUpdate {
			if filterUpdate.ID == filterFeed.FilterId {
				filtersUpdate[fsKey].FeedIds = append(filtersUpdate[fsKey].FeedIds, filterFeed.FeedId)
			}
		}
	}
	return filtersUpdate
}
