package item_model

import "github.com/jinzhu/gorm"

const TagItemTable = "later_item"

type TagItem struct {
	gorm.Model

	UserItemId uint
	TagId      uint `gorm:"column:later_id"`
	Unread     int
}

func (TagItem) TableName() string {
	return TagItemTable
}

func (item TagItem) IsUnread() bool {
	if (item.Unread > 0) {
		return true
	}

	return false
}
