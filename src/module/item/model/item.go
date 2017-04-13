package item_model

import (
	"time"
	"github.com/jinzhu/gorm"
)

const ItemTable = "item"

type Item struct {
	gorm.Model

	Title       string
	Link        string
	Content     string
	ContentHash string `gorm:"column:content_hash"`
	Author      string
	FeedId      uint   `gorm:"column:feed_id"`
	PublishedAt time.Time
}

func (Item) TableName() string {
	return ItemTable
}

