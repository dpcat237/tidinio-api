package item_model

import (
	"github.com/jinzhu/gorm"
	"time"
	"github.com/cstockton/go-conv"
)

const UserItemTable = "user_item"

type UserItem struct {
	gorm.Model

	ItemId uint
	UserId uint
	Unread int
	Stared int
	Shared int
}

func (UserItem) TableName() string {
	return UserItemTable
}

func FromUserItemsSync(items []UserItemSync) []UserItem {
	result := []UserItem{}
	for _, item := range items {
		uItem := UserItem{}
		uItem.ID = item.ID
		uItem.Unread, _ = conv.Int(item.Unread)
		uItem.Stared, _ = conv.Int(item.Stared)
		result = append(result, uItem)
	}

	return result
}

type UserItemSync struct {
	ID        uint      `json:"article_id"`
	FeedId    uint      `json:"feed_id"`
	Title     string    `json:"title"`
	Link      string    `json:"link"`
	Content   string    `json:"content"`
	Stared    bool      `json:"is_stared"`
	Unread    bool      `json:"is_unread"`
	PublishedAt time.Time `json:"date_add"`
}

type SharedItem struct {
	TagId uint   `json:"tag_id"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

func (item UserItem) IsUnread() bool {
	if (item.Unread > 0) {
		return true
	}

	return false
}

func MergeToUserItemSync(item Item, userItem UserItem) UserItemSync {
	result := UserItemSync{}
	result.ID = userItem.ID
	result.FeedId = item.FeedId
	result.Title = item.Title
	result.Link = item.Link
	result.Content = item.Content
	result.Stared, _ = conv.Bool(userItem.Stared)
	result.Unread, _ = conv.Bool(userItem.Unread)
	result.PublishedAt = item.PublishedAt

	return result
}

func MergeUserItems(collection1 []UserItem, collection2 []UserItem) []UserItem {
	for _, value := range collection2 {
		collection1 = append(collection1, value)
	}

	return collection1
}

func ReadToUserItemSync(userItem UserItem) UserItemSync {
	result := UserItemSync{}
	result.ID = userItem.ID
	result.FeedId = 0
	result.Stared = false
	result.Unread = false

	return result
}

func (userItem *UserItem) SetUnread() {
	userItem.Unread = 1
}
