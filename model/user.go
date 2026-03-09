package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        int64          `gorm:"column:id;type:bigint(20);primary_key" json:"id,string"`
	Phone     string         `gorm:"column:phone;type:varchar(255);comment:手机号;NOT NULL" json:"phone"`
	Nickname  string         `gorm:"column:nickname;type:varchar(255);comment:昵称;NOT NULL" json:"nickname"`
	Avatar    string         `gorm:"column:avatar;type:varchar(255);comment:头像" json:"avatar"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}

func (m *User) TableName() string {
	return "user"
}
