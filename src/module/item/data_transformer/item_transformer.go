package item_transformer

import (
	"github.com/cstockton/go-conv"
	"github.com/tidinio/src/module/item/model"
)

func FromUserItemsSync(items []item_model.UserItemSync) []item_model.UserItem {
	result := []item_model.UserItem{}
	for _, item := range items {
		uItem := item_model.UserItem{}
		uItem.ID = item.ID
		uItem.Unread, _ = conv.Int(item.Unread)
		uItem.Stared, _ = conv.Int(item.Stared)
		result = append(result, uItem)
	}

	return result
}

func ToTagItemSync(items []item_model.TagItemSyncDB) []item_model.TagItemSync {
	result := []item_model.TagItemSync{}
	orderedTags := moveTagItemsUnderUserItemId(items)
	for userItemId, tags := range orderedTags {
		item := item_model.TagItemSync{}
		item.ArticleId = userItemId
		item.Stared = tags[0].Stared
		item.Tags = getTagsIds(tags)
		result = append(result, item)
	}

	return result
}

func getTagsIds(collection []item_model.TagItemSyncDB) []uint {
	result := []uint{}
	for _, item := range collection {
		result = append(result, item.TagId)
	}

	return result
}

func moveTagItemsUnderUserItemId(collection []item_model.TagItemSyncDB) map[uint][]item_model.TagItemSyncDB {
	joined := make(map[uint][]item_model.TagItemSyncDB)
	for _, tagItem := range collection {
		joined[tagItem.ArticleId] = append(joined[tagItem.ArticleId], tagItem)
	}

	return joined
}
