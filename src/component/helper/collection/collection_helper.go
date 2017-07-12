package collection_helper

import (
	"fmt"

	"github.com/tidinio/src/module/item/model"
)

func CheckIdExistInCollection(id uint, collection []uint) bool {
	for _, key := range collection {
		if key == id {
			return true
		}
	}

	return false
}

func ConvertIntToStringSlice(collection []uint) []string {
	result := make([]string, len(collection))
	for key, item := range collection {
		result[key] = fmt.Sprint(item)
	}

	return result
}

func GetIdsFromUserItemCollectionStr(collection []item_model.UserItem) []string {
	ids := make([]string, len(collection))
	for key, item := range collection {
		ids[key] = fmt.Sprint(item.ID)
	}

	return ids
}

func GetItemIdsFromUserItemCollectionStr(collection []item_model.UserItem) []string {
	ids := make([]string, len(collection))
	for key, item := range collection {
		ids[key] = fmt.Sprint(item.ItemId)
	}

	return ids
}

func GetUserItemIdsFromTagItemCollectionStr(collection []item_model.TagItem) []string {
	ids := make([]string, len(collection))
	for key, item := range collection {
		ids[key] = fmt.Sprint(item.UserItemId)
	}

	return ids
}

func GetUserItemIdsFromTagItemListCollectionStr(collection []item_model.TagItemList) []string {
	ids := make([]string, len(collection))
	for key, item := range collection {
		ids[key] = fmt.Sprint(item.ArticleId)
	}

	return ids
}

func MoveTagItemsUnderUserItemId(collection []item_model.TagItem) map[uint][]item_model.TagItem {
	joined := make(map[uint][]item_model.TagItem)
	for _, tagItem := range collection {
		joined[tagItem.UserItemId] = append(joined[tagItem.UserItemId], tagItem)
	}

	return joined
}

func MoveTagItemsListUnderUserItemId(collection []item_model.TagItemList) map[uint]item_model.TagItemList {
	joined := make(map[uint]item_model.TagItemList)
	for _, tagItem := range collection {
		joined[tagItem.ArticleId] = tagItem
	}

	return joined
}
