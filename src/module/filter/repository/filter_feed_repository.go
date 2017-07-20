package filter_repository

import (
	"fmt"

	"github.com/tidinio/src/component/repository"
	"github.com/tidinio/src/module/filter/model"
)

const filterFeedTable = filter_model.FilterFeedTable

func DeleteFilterFeeds(filterId uint) {
	query := fmt.Sprintf("DELETE FROM "+filterFeedTable+" WHERE filter_id = '%d'", filterId)
	app_repository.Conn.Exec(query)
}

func GetFilterFeedsByUserId(userId uint) []filter_model.FilterFeed {
	filterFeeds := []filter_model.FilterFeed{}
	app_repository.Conn.Table(filterFeedTable).
		Joins("inner join " + filterTable + " on " + filterFeedTable + ".filter_id = " + filterTable + ".id").
		Where(filterTable+".user_id = ? and "+filterTable+".deleted_at IS NULL", userId).
		Scan(&filterFeeds)
	return filterFeeds
}

func SaveFilterFeed(filterFeed *filter_model.FilterFeed) {
	if app_repository.Conn.NewRecord(filterFeed) {
		app_repository.Conn.Create(&filterFeed)
	} else {
		app_repository.Conn.Save(&filterFeed)
	}
}
