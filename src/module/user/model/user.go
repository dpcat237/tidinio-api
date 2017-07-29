package user_model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	UserBasic

	Enabled      int
	Registered   int
	Subscribed   int
	PreferenceId uint
}

type UserBasic struct {
	gorm.Model

	Email    string `gorm:"size:255;not null;unique"`
	Password string `gorm:"size:255"`
	Preview  int
}

type UserFeedback struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Email string `gorm:"size:255"`
	Title string `gorm:"size:255"`
	Text  string `gorm:"size:5000"`
}

type UserFeedbackSync struct {
	Email string `json:"email"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type UserSync struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (UserBasic) TableName() string {
	return "user"
}

func (UserFeedback) TableName() string {
	return "feedback"
}
