package example

import (
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	UserID  uint `gorm:"not null;uniqueIndex"` // 外键字段，确保一对一关系
	Age     int
	Address string `gorm:"type:varchar(200)"`
}
