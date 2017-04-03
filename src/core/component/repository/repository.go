package app_repository

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"time"
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

func (repo Repository) Close()  {
	repo.DB.Close()
}

func GetDateNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func InitConnection() Repository {
	repo := Repository{}
	db, err := gorm.Open("mysql", "tidinio:pwd@tcp(dbcontainer:3306)/tidinio?charset=utf8&parseTime=True")
	if err != nil {
		//log.Println("err happened", err)
		fmt.Println("err happened", err)
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
