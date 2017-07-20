package filter_model

const FilterFeedTable = "filter_feed"

type FilterFeed struct {
	ID uint `gorm:"primary_key"`

	FeedId   uint `gorm:"column:feed_id"`
	FilterId uint `gorm:"column:filter_id"`
}

func (FilterFeed) TableName() string {
	return FilterFeedTable
}
