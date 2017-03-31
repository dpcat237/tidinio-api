package feed_model

import (
	"github.com/jinzhu/gorm"
)

const FeedHistoryTable = "feed_history"

type FeedHistory struct {
	gorm.Model

	FeedId       uint `gorm:"column:feed_id"`
	CountWaiting int  `gorm:"column:count_waiting"`
	Finished     int
}

func (FeedHistory) TableName() string {
	return FeedHistoryTable
}

func (feedHistory *FeedHistory) IncreaseCounter() {
	feedHistory.CountWaiting++
}

func (feedHistory FeedHistory) IsFinished() bool {
	return (feedHistory.Finished == 1)
}

func (feedHistory *FeedHistory) SetFinished() {
	feedHistory.Finished = 1
}
