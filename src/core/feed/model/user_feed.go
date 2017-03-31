package feed_model

import (
	"time"
)

const UserFeedTable = "user_feed"

type UserFeed struct {
	ID        uint `gorm:"primary_key"`

	FeedId    uint
	UserId    uint
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserFeed) TableName() string {
	return UserFeedTable
}
