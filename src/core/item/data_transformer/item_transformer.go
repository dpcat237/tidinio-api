package item_transformer

import (
	"github.com/tidinio/src/core/component/repository"
	"github.com/tidinio/src/core/item/model"
)

func FromUserItemsSync(items []item_model.UserItemSync) []item_model.UserItem {
	result := []item_model.UserItem{}
	for _, item := range items {
		uItem := item_model.UserItem{}
		uItem.ID = item.ID
		uItem.Unread = common_repository.BoolToInt(item.Unread)
		uItem.Stared = common_repository.BoolToInt(item.Stared)
		result = append(result, uItem)
	}

	return result
}
