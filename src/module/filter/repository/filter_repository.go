package filter_repository

import (
	"github.com/tidinio/src/component/repository"
	"github.com/tidinio/src/module/filter/model"
)

const filterTable = filter_model.FilterTable

func DeleteFilter(filter filter_model.Filter) {
	app_repository.Conn.Delete(filter)
}

func GetUserFilters(userId uint) []filter_model.Filter {
	filters := []filter_model.Filter{}
	app_repository.Conn.Table(filterTable).Where("user_id = ? and deleted_at IS NULL", userId).Scan(&filters)
	return filters
}

func SaveFilter(filter *filter_model.Filter) {
	if app_repository.Conn.NewRecord(filter) {
		app_repository.Conn.Create(&filter)
	} else {
		app_repository.Conn.Save(&filter)
	}
}
