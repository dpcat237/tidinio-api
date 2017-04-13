package app_repository

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"time"
	"github.com/tidinio/src/core/component/logger"
)

const (
	db_host = "dbcontainer"
	db_user = "tidinio"
	db_password = "pwd"
	db_name = "tidinio"
)

type Repository struct {
	DB *gorm.DB
}

func BoolToInt(value bool) int {
	if (value) {
		return 1
	}

	return 0
}

func (repo Repository) Close() {
	repo.DB.Close()
}

func GetDateNow() *time.Time {
	now := time.Now()

	return &now
}

func GetDateNowFormatted() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func InitConnection() Repository {
	repo := Repository{}
	db, err := gorm.Open("mysql", db_user + ":" + db_password + "@tcp(" + db_host + ":3306)/" + db_name + "?charset=utf8&parseTime=True")
	if err != nil {
		app_logger.Error(err.Error())
	}
	//db.DB().SetMaxIdleConns(1)
	db.LogMode(true)
	repo.DB = db

	return repo
}

func IntToBool(value int) bool {
	if (value > 0) {
		return true
	}

	return false
}

func UintToString(value uint) string {
	return fmt.Sprint(value)
}
