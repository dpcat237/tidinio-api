package feed_model

import (
	"time"
	"github.com/jinzhu/gorm"
)

const UserFeedTable = "user_feed"

type UserFeed struct {
	gorm.Model

	FeedId    uint
	UserId    uint
	Title     string
}

type UserFeedSync struct {
	ID        uint `json:"feed_id"`

	FeedId    uint
	UserId    uint
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserFeed) TableName() string {
	return UserFeedTable
}
