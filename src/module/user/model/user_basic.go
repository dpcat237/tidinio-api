package user_model

import (
	"github.com/jinzhu/gorm"
)

type UserBasic struct {
	gorm.Model

	Email    string `gorm:"size:255;not null;unique"`
	Password string `gorm:"size:255"`
}

func CreateUserBasic(email string, password string) UserBasic {
	user := UserBasic{}
	user.Email = email
	user.Password = password

	return user
}

func (UserBasic) TableName() string {
	return "user"
}
