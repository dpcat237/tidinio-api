package app_repository

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
	"github.com/tidinio/src/component/configuration"
	"github.com/tidinio/src/component/logger"
)

var Conn = gorm.DB{}

func Close() {
	Conn.Close()
}

func GetDateNow() *time.Time {
	now := time.Now()

	return &now
}

func GetDateNowFormatted() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func InitConnection() {
	dbUser, _ := app_conf.Data.String("database.user")
	dbPassword, _ := app_conf.Data.String("database.password")
	dbHost, _ := app_conf.Data.String("database.host")
	dbName, _ := app_conf.Data.String("database.name")
	db, err := gorm.Open("mysql", dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":3306)/" + dbName + "?charset=utf8&parseTime=True")
	if err != nil {
		app_logger.Error(err.Error())
	}
	//db.DB().SetMaxIdleConns(1)
	db.LogMode(true)

	Conn = *db
}
