package feed_model

import (
	"time"
	"github.com/jinzhu/gorm"
)

const UserFeedTable = "user_feed"

type UserFeed struct {
	gorm.Model

	FeedId uint
	UserId uint
	Title  string
}

type UserFeedSync struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"date_up"`
}

func (UserFeed) TableName() string {
	return UserFeedTable
}

func ToUserFeedsSync(userFeeds []UserFeed) []UserFeedSync {
	userFeedsSync := []UserFeedSync{}
	for _, userFeed := range userFeeds {
		userFeedsSync = append(userFeedsSync, ToUserFeedSync(userFeed))
	}

	return userFeedsSync
}

func ToUserFeedSync(userFeed UserFeed) UserFeedSync {
	userFeedSync := UserFeedSync{}
	userFeedSync.ID = userFeed.ID
	userFeedSync.Title = userFeed.Title
	userFeedSync.UpdatedAt = userFeed.UpdatedAt

	return userFeedSync
}
