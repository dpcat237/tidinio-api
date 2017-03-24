package helper_collection

import (
	"github.com/tidinio/src/core/item/model"
	"fmt"
)

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
