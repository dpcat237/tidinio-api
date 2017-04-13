package app

import (
	"github.com/tidinio/src/component/configuration"
	"github.com/tidinio/src/component/repository"
)

func InitializeRequiredData() {
	app_conf.LoadConfiguration()
	app_repository.InitConnection()
}
