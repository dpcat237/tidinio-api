package item_model

import "github.com/jinzhu/gorm"

const ItemTable = "item"

type Item struct {
	gorm.Model

	Title string
	Link string
	Content string
	ContentHash  string
	Author  string
	FeedId uint
}

