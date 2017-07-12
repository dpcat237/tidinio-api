package tag_model

import (
	"github.com/jinzhu/gorm"
	"time"
)

const TagTable = "tag"

type Tag struct {
	gorm.Model

	Name    string
	UserId  uint `gorm:"column:user_id"`
	Enabled int
	Basic   int `gorm:"column:is_basic"`
}

type TagSync struct {
	ID        uint      `json:"tag_id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"date_up"`
}

func (tag Tag) ConvertToTagSync() TagSync {
	tagSync := TagSync{}
	tagSync.ID = tag.ID
	tagSync.Name = tag.Name
	tagSync.UpdatedAt = tag.UpdatedAt
	return tagSync
}

func ConvertToTagsSync(tags []Tag) []TagSync {
	tagsSync := []TagSync{}
	for _, tag := range tags {
		tagsSync = append(tagsSync, tag.ConvertToTagSync())
	}
	return tagsSync
}

func (Tag) TableName() string {
	return TagTable
}
