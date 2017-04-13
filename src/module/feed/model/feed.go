package feed_model

import (
	"github.com/jinzhu/gorm"
	"time"
)

const FeedTable = "feed"

type Feed struct {
	gorm.Model

	Title        string
	Url          string
	Website      string
	Language     string
	Favicon      string
	Enabled      int
	Crawling     int `gorm:"column:is_crawling"`
	SyncInterval int `gorm:"column:sync_interval"`
	SyncAt       time.Time
}

func (Feed) TableName() string {
	return FeedTable
}

func (feed *Feed) Enable() {
	feed.Enabled = 1
}

func (feed Feed) IsEnabled() bool {
	return (feed.Enabled == 1)
}

func (feed *Feed) UpdateDateSync() {
	feed.SyncAt = time.Now()
}
