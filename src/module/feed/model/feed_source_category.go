package feed_model

import "github.com/jinzhu/gorm"

const FeedSourceCategoryTable = "feed_source_category"

type FeedSourceCategory struct {
	gorm.Model

	Name        string
	ImageName   string
	ImageFile   string
	FeedSources []FeedSource `gorm:"many2many:feed_sources_categories;"`
}

type FeedSourceCategorySync struct {
	ID          uint         		 `json:"id"`
	Name        string       		 `json:"name"`
	ImageName   string       		 `json:"image_name"`
	ImageFile   string       		 `json:"image_file"`
	FeedSources []FeedSourceSync `json:"feed_sources"`
}

func (FeedSourceCategory) TableName() string {
	return FeedSourceCategoryTable
}

func ToFeedSourceCategorySyncs(categories []FeedSourceCategory) []FeedSourceCategorySync {
	categoriesSync := []FeedSourceCategorySync{}
	for _, category := range categories {
		categoriesSync = append(categoriesSync, toFeedSourceCategorySync(category))
	}

	return categoriesSync
}

func toFeedSourceCategorySync(category FeedSourceCategory) FeedSourceCategorySync {
	categorySync := FeedSourceCategorySync{}
	categorySync.ID = category.ID
	categorySync.Name = category.Name
	categorySync.ImageName = category.ImageName
	categorySync.ImageFile = category.ImageFile
	categorySync.FeedSources = toFeedSourceSyncs(category.FeedSources)

	return categorySync
}
