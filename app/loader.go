package app

import (
	"github.com/tidinio/src/core/component/configuration"
	"github.com/tidinio/src/core/component/repository"
)

func InitializeRequiredData() {
	app_conf.LoadConfiguration()
	app_repository.InitConnection()
}
