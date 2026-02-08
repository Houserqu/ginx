package model

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	Id             int64          `gorm:"column:id;type:bigint(20);primary_key" json:"id,string"`
	Name           string         `gorm:"column:name;type:varchar(255);comment:团队名称;NOT NULL" json:"name"`
	Creator        int64          `gorm:"column:creator;type:bigint(20);comment:创建人" json:"creator,string"`
	Avatar         string         `gorm:"column:avatar;type:varchar(255);comment:头像" json:"avatar"`
	Compute        int            `gorm:"column:compute;type:int(11);default:0;comment:计算值" json:"compute"`
	DataStatistics int            `gorm:"column:data_statistics;type:int(11);default:0;comment:数据统计开关，0 关闭，1 开启" json:"data_statistics"`
	ApiKey         string         `gorm:"column:api_key;type:varchar(255);comment:API密钥" json:"-"`
	CreatedAt      time.Time      `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}

func (m *Team) TableName() string {
	return "team"
}
