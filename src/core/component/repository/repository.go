package common_repository

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
)

func InitConnection() (db *gorm.DB) {
	db, err := gorm.Open("mysql", "tidinio:pwd@tcp(dbcontainer:3306)/tidinio?charset=utf8&parseTime=True")
	if err != nil {
		//log.Println("err happened", err)
		fmt.Println("err happened", err)
	}
	//defer db.Close()
	//db.DB().SetMaxIdleConns(1)
	db.LogMode(true)

	return db
}

func InitConnectionP() (db *gorm.DB) {
	db, err := gorm.Open("postgres", "host=postgres user=tidinio dbname=tidinio sslmode=disable password=pwd")
	if err != nil {
		//log.Println("err happened", err)
		fmt.Println("err happened", err)
	}
	//defer db.Close()
	//db.DB().SetMaxIdleConns(1)
	db.LogMode(true)

	return db
}

func BoolToInt(value bool) int {
	if (value) {
		return 1
	}

	return 0
}

func IntToBool(value int) bool {
	if (value > 0) {
		return true
	}

	return false
}
