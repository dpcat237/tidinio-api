package helper_collection

import (
	"github.com/tidinio/src/core/item/model"
	"fmt"
)

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
