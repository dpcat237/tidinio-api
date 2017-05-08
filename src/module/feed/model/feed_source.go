package feed_model

import (
	"github.com/jinzhu/gorm"
)

const FeedSourceTable = "feed_source"

type FeedSource struct {
	gorm.Model

	Name    string
	Web     string
	FeedUrl string
}

type FeedSourceSync struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Web     string `json:"web"`
	FeedUrl string `json:"feed_url"`
}

func (FeedSource) TableName() string {
	return FeedSourceTable
}

func toFeedSourceSyncs(sources []FeedSource) []FeedSourceSync {
	sourcesSync := []FeedSourceSync{}
	for _, source := range sources {
		sourcesSync = append(sourcesSync, toFeedSourceSync(source))
	}

	return sourcesSync
}

func toFeedSourceSync(source FeedSource) FeedSourceSync {
	sourceSync := FeedSourceSync{}
	sourceSync.ID = source.ID
	sourceSync.Name = source.Name
	sourceSync.Web = source.Web
	sourceSync.FeedUrl = source.FeedUrl

	return sourceSync
}
