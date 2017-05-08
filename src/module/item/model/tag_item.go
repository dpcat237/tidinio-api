package item_model

import (
	"github.com/jinzhu/gorm"
	"time"
)

const TagItemTable = "tag_item"

type TagItem struct {
	gorm.Model

	UserItemId uint
	TagId      uint `gorm:"column:tag_id"`
	Unread     int
}

type TagItemList struct {
	ArticleId   uint      `json:"article_id" gorm:"column:user_item_id"`
	FeedId      uint      `json:"feed_id"  gorm:"column:feed_id"`
	Language    string    `json:"language"`
	Stared      bool      `json:"is_stared"`
	Link        string    `json:"link"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	PublishedAt time.Time `json:"date_add"`
	Tags        []uint    `json:"tags"`
}

type TagItemSync struct {
	ArticleId uint      `json:"article_id" gorm:"column:user_item_id"`
	Stared    bool      `json:"is_stared"`
	Tags      []uint    `json:"tags"`
}

type TagItemSyncDB struct {
	ArticleId uint      `gorm:"column:user_item_id"`
	Stared    bool      `json:"is_stared"`
	TagId     uint      `gorm:"column:tag_id"`
}

func (TagItem) TableName() string {
	return TagItemTable
}

func (item TagItem) IsUnread() bool {
	if item.Unread > 0 {
		return true
	}

	return false
}

func AddTagItemsListContent(tagItems []TagItemList, tagItemsContent []TagItemList) []TagItemList {
	for _, tagItem := range tagItems {
		for idx, tagItemContent := range tagItemsContent {
			if tagItem.ArticleId == tagItemContent.ArticleId {
				tagItemsContent[idx].Tags = tagItem.Tags
			}
		}
	}

	return tagItemsContent
}

func (item TagItemList) HasTag(tagId uint) bool {
	for _, tag := range item.Tags {
		if tag == tagId {
			return true
		}
	}

	return false
}

func JoinTagsByUserItem(tagItems []TagItem) map[uint][]uint {
	joined := make(map[uint][]uint)
	for _, tagItem := range tagItems {
		joined[tagItem.UserItemId] = append(joined[tagItem.UserItemId], tagItem.TagId)
	}

	return joined
}

func MergeTagItemsList(collection1 []TagItemList, collection2 []TagItemList) []TagItemList {
	for _, value := range collection2 {
		collection1 = append(collection1, value)
	}

	return collection1
}

func MergeToTagItemList(userItemId uint, tagsIds []uint) TagItemList {
	item := TagItemList{}
	item.ArticleId = userItemId
	item.Tags = tagsIds

	return item
}

func MergeToTagItemsList(tagItems []TagItem, relatedTags map[uint][]uint) []TagItemList {
	results := []TagItemList{}
	for _, tagItem := range tagItems {
		for userItemId, tags := range relatedTags {
			if (tagItem.UserItemId == userItemId) {
				results = append(results, MergeToTagItemList(tagItem.UserItemId, tags))
			}
		}
	}

	return results
}

func ToTagItemSync(items []TagItemSyncDB) []TagItemSync {
	result := []TagItemSync{}
	orderedTags := moveTagItemsUnderUserItemId(items)
	for userItemId, tags := range orderedTags {
		item := TagItemSync{}
		item.ArticleId = userItemId
		item.Stared = tags[0].Stared
		item.Tags = getTagsIds(tags)
		result = append(result, item)
	}

	return result
}

func getTagsIds(collection []TagItemSyncDB) []uint {
	result := []uint{}
	for _, item := range collection {
		result = append(result, item.TagId)
	}

	return result
}

func moveTagItemsUnderUserItemId(collection []TagItemSyncDB) map[uint][]TagItemSyncDB {
	joined := make(map[uint][]TagItemSyncDB)
	for _, tagItem := range collection {
		joined[tagItem.ArticleId] = append(joined[tagItem.ArticleId], tagItem)
	}

	return joined
}
